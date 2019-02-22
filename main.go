package main

import (
	"encoding/json"
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"
	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/jhunt/go-table"

	"github.com/jhunt/osb/api"
)

var Version = "(development version)"

var opt struct {
	Help  bool `cli:"-h, --help"`
	Trace bool `cli:"-T, --trace" env:"OSB_TRACE"`

	Version bool `cli:"-v, --version"`

	Data string `cli:"--data" env:"OSB_DATA"`

	Endpoint   string `cli:"-e, --endpoint" env:"OSB_URL"`
	Username   string `cli:"-U, --username" env:"OSB_USERNAME"`
	Password   string `cli:"-P, --password" env:"OSB_PASSWORD"`
	SkipVerify bool   `cli:"-k, --skip-verify" env:"OSB_SKIP_VERIFY"`
	Timeout    int    `cli:"-t, --timeout" env:"OSB_TIMEOUT"`

	JSON bool `cli:"--json"`

	List struct{} `cli:"list, ls"`
	Env  struct{} `cli:"env"`

	Catalog struct{} `cli:"catalog"`

	Provision struct {
		ID string `cli:"-i, --instance, --id"`
	} `cli:"provision, prov, create"`

	Bind struct {
		Service string `cli:"-s, --service"`
		Plan    string `cli:"-p, --plan"`
		ID      string `cli:"-i, --binding, --id"`
	} `cli:"bind"`

	Unbind struct {
		Service string `cli:"-s, --service"`
		Plan    string `cli:"-p, --plan"`
		ID      string `cli:"-i, --binding, --id"`
	} `cli:"unbind"`

	Deprovision struct {
		Service string `cli:"-s, --service"`
		Plan    string `cli:"-p, --plan"`
	} `cli:"deprovision, deprov, rm"`
}

func bail(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", e)
		os.Exit(1)
	}
}

func main() {
	opt.Timeout = 5
	env.Override(&opt)
	command, args, err := cli.Parse(&opt)
	bail(err)

	if command == "" && len(args) == 0 && opt.Help {
		fmt.Printf("USAGE: @G{%s} [@W{options}] <@C{command}> [@W{options}]\n\n", os.Args[0])
		fmt.Printf("Options:\n\n")
		fmt.Printf("  -h, --help         Show the help screen.\n")
		fmt.Printf("\n")
		fmt.Printf("  -T, --trace        Trace HTTP requests and responses as they happen.\n")
		fmt.Printf("                     Can also be enabled by setting @W{OSB_TRACE=yes}.\n")
		fmt.Printf("\n")
		fmt.Printf("  --data             Path to the OSB data file, for storing instance and\n")
		fmt.Printf("                     binding information required by future bind, unbind,\n")
		fmt.Printf("                     and deprovision requests.\n")
		fmt.Printf("\n")
		fmt.Printf("  -e, --endpoint     The URL to the backend service broker to interact with.\n")
		fmt.Printf("                     Can also be specified via the @W{OSB_URL} variable.\n")
		fmt.Printf("\n")
		fmt.Printf("  -U, --username     The username for service broker HTTP Basic Auth.\n")
		fmt.Printf("                     Can also be specified via the @W{OSB_USERNAME} variable.\n")
		fmt.Printf("\n")
		fmt.Printf("  -P, --password     The password for service broker HTTP Basic Auth.\n")
		fmt.Printf("                     Can also be specified via the @W{OSB_PASSWORD} variable.\n")
		fmt.Printf("\n")
		fmt.Printf("  -k, --skip-verify  Do not validate X.509 TLS certificates.\n")
		fmt.Printf("                     Can also be specified via @W{OSB_SKIP_VERIFY}.\n")
		fmt.Printf("\n")
		fmt.Printf("  -t, --timeout      Timeout (in seconds) for HTTP reuests.\n")
		fmt.Printf("                     Can also be specified via @W{OSB_TIMEOUT}.\n")
		fmt.Printf("\n")
		fmt.Printf("  --json             Emit JSON responses, and nothing else.\n")
		fmt.Printf("                     Useful for scripting!\n")
		fmt.Printf("\n")
		fmt.Printf("Commands:\n\n")
		fmt.Printf("  list           List known instance and binding details, from ~/.osbrc.\n")
		fmt.Printf("  env            Dump the environment variables that `osb` cares about.\n")
		fmt.Printf("  catalog        Retrieve the service catalog from the service broker.\n")
		fmt.Printf("\n")
		fmt.Printf("  provision      Provision a new instance of a service/plan.\n")
		fmt.Printf("  deprovision    Remove a provsioned instance.\n")
		fmt.Printf("\n")
		fmt.Printf("  bind           Bind a provisioned instance, to get credentials.\n")
		fmt.Printf("  unbind         Unbind an instance, releasing bound credentials.\n")
		fmt.Printf("\n")
		os.Exit(0)
	}

	if command == "" && len(args) == 0 && opt.Version {
		fmt.Printf("osb %s\n", Version)
		os.Exit(0)
	}

	if command == "" && len(args) > 0 {
		fmt.Fprintf(os.Stderr, "@R{Unrecognized command '%s'}\n", args[0])
		os.Exit(1)
	}

	c := &api.Client{
		URL:        opt.Endpoint,
		Username:   opt.Username,
		Password:   opt.Password,
		SkipVerify: opt.SkipVerify,
		Timeout:    opt.Timeout,
		Trace:      opt.Trace,
	}

	store, err := api.ReadStore(opt.Data)
	bail(err)

	switch command {
	default:
		bail(fmt.Errorf("%s: not implemented", command))

	case "list":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s}\n\n", os.Args[0], command)
			os.Exit(0)
		}

		if opt.JSON {
			jsonify(store)
			os.Exit(0)
		}

		t := table.NewTable("Broker", "Instance", "Service", "Plan", "Binding", "Credentials")
		for _, broker := range store.Data {
			bname := broker.Broker
			for _, instance := range broker.Instances {
				if instance.Bindings == nil || len(instance.Bindings) == 0 {
					t.Row(nil, bname, instance.ID, instance.ServiceID, instance.PlanID, "-", "-")
					bname = ""

				} else {
					inst := instance.ID
					service := instance.ServiceID
					plan := instance.PlanID
					for _, binding := range instance.Bindings {
						b, err := json.MarshalIndent(binding.Credentials, "", "  ")
						if err != nil {
							t.Row(nil, bname, inst, service, plan, binding.ID, fmt.Sprintf("error: %s", err))
						} else {
							t.Row(nil, bname, inst, service, plan, binding.ID, string(b))
						}
						inst = ""
						service = ""
						plan = ""
					}
				}
			}
		}
		t.Output(os.Stdout)
		os.Exit(0)

	case "env":
		e := struct {
			Trace      bool   `json:"OSB_TRACE"`
			Data       string `json:"OSB_DATA"`
			URL        string `json:"OSB_URL"`
			Username   string `json:"OSB_USERNAME"`
			Password   string `json:"OSB_PASSWORD"`
			SkipVerify bool   `json:"OSB_SKIP_VERIFY"`
			Timeout    int    `json:"OSB_TIMEOUT"`
		}{
			Trace:      opt.Trace,
			Data:       opt.Data,
			URL:        opt.Endpoint,
			Username:   opt.Username,
			Password:   opt.Password,
			SkipVerify: opt.SkipVerify,
			Timeout:    opt.Timeout,
		}

		if opt.JSON {
			jsonify(e)
			os.Exit(0)
		}

		booly := func(tf bool) string {
			if tf {
				return "yes"
			} else {
				return "no"
			}
		}

		fmt.Printf("export OSB_URL=\"%s\"\n", e.URL)
		fmt.Printf("export OSB_USERNAME=\"%s\"\n", e.Username)
		fmt.Printf("export OSB_PASSWORD=\"%s\"\n", e.Password)
		fmt.Printf("export OSB_TIMEOUT=%d\n", e.Timeout)
		fmt.Printf("export OSB_DATA=\"%s\"\n", e.Data)
		fmt.Printf("export OSB_TRACE=%s\n", booly(e.Trace))
		fmt.Printf("export OSB_SKIP_VERIFY=%s\n", booly(e.SkipVerify))

	case "catalog":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s}\n\n", os.Args[0], command)
			os.Exit(0)
		}

		connecting()
		catalog, err := c.GetCatalog()
		bail(err)

		if opt.JSON {
			jsonify(catalog)
			os.Exit(0)
		}

		t := table.NewTable("Service", "(ID)", "Plans", "(IDs)", "Tags")
		for _, s := range catalog.Services {

			plans := ""
			ids := ""
			for _, p := range s.Plans {
				plans += fmt.Sprintf("%s\n", p.Name)
				ids += fmt.Sprintf("%s\n", p.ID)
			}
			if plans == "" {
				plans = "(none)"
			}

			tags := ""
			for _, t := range s.Tags {
				tags += fmt.Sprintf("%s\n", t)
			}
			if tags == "" {
				tags = "(none)"
			}

			t.Row(nil, s.Name, s.ID, plans, ids, tags)
			t.Row(nil, "", "", "", "", "")
		}
		t.Output(os.Stdout)
		os.Exit(0)

	case "provision":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s} [--id ID] @M{SERVICE}/@M{PLAN}\n\n", os.Args[0], command)
			fmt.Printf("Options:\n\n")
			fmt.Printf("  -i, --id       The ID to use for the newly-provisioned service instance.\n")
			fmt.Printf("                 If not specified, will be a random UUID.\n")
			fmt.Printf("\n")
			os.Exit(0)
		}

		provisioning(args)
		l := strings.SplitN(args[0], "/", 2)
		if len(l) != 2 {
			fmt.Fprintf(os.Stderr, "USAGE: @G{%s} [@W{options}] @C{%s} [--id ID] @M{SERVICE}/@M{PLAN}\n", os.Args[0], command)
			fmt.Fprintf(os.Stderr, "@R{(missing the /plan bits...)}\n\n")
			os.Exit(1)
		}

		catalog, err := c.GetCatalog()
		bail(err)

		service, plan, err := catalog.FindPlan(l[0], l[1])
		bail(err)

		stat, err := c.Provision(opt.Provision.ID, api.ProvisionSpec{
			ServiceID: service,
			PlanID:    plan,
		})
		bail(err)

		store.AddInstance(c.URL, stat.InstanceID, service, plan)
		if err := store.Write(opt.Data); err != nil {
			fmt.Fprintf(os.Stderr, "@Y{!!! %s}\n", err)
		}

		if opt.JSON {
			jsonify(stat)
			os.Exit(0)
		}

		fmt.Printf("instance: @G{%s}\n", stat.InstanceID)
		fmt.Printf("status:   @M{%s}\n", stat.Status)
		if stat.DashboardURL != "" {
			fmt.Printf("dashboard: @C{%s}\n", stat.DashboardURL)
		}
		if stat.Operation != "" {
			fmt.Printf("operation: @C{%s}\n", stat.Operation)
		}
		os.Exit(0)

	case "bind":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s} [@W{options}] @M{INSTANCE}\n\n", os.Args[0], command)
			fmt.Printf("Options:\n\n")
			fmt.Printf("  -s, --service  The name or ID of the service that the instance\n")
			fmt.Printf("                 was provisioned from.  This is required if the\n")
			fmt.Printf("                 instance details are not found in ~/.osbrc.\n")
			fmt.Printf("\n")
			fmt.Printf("  -p, --plan     The name or ID of the plan that the instance\n")
			fmt.Printf("                 was provisioned from.  This is required if the\n")
			fmt.Printf("                 instance details are not found in ~/.osbrc.\n")
			fmt.Printf("\n")
			fmt.Printf("  -i, --id       The ID to use for the new service instance binding.\n")
			fmt.Printf("                 If not specified, will be a random UUID.\n")
			fmt.Printf("\n")
			os.Exit(0)
		}

		binding(args)

		service, plan, _ := store.GetInstanceDetails(c.URL, args[0])
		if service == "" || plan == "" {
			catalog, err := c.GetCatalog()
			bail(err)

			if service == "" {
				service = opt.Bind.Service
				if service == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", args[0])
					fmt.Fprintf(os.Stderr, "You must specify the --service flag to the bind operation.\n")
					os.Exit(1)
				}
			}
			if plan == "" {
				plan = opt.Bind.Plan
				if plan == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", args[0])
					fmt.Fprintf(os.Stderr, "You must specify the --plan flag to the bind operation.\n")
					os.Exit(1)
				}
			}

			s, p, err := catalog.FindPlan(service, plan)
			bail(err)
			service = s
			plan = p
		}

		stat, err := c.Bind(api.BindSpec{
			InstanceID: args[0],
			BindingID:  opt.Bind.ID,
			ServiceID:  service,
			PlanID:     plan,
		})
		bail(err)

		store.AddBinding(c.URL, stat.InstanceID, stat.BindingID, stat.Credentials)
		if err := store.Write(opt.Data); err != nil {
			fmt.Fprintf(os.Stderr, "@Y{!!! %s}\n", err)
		}

		if opt.JSON {
			jsonify(stat)
			os.Exit(0)
		}

		fmt.Printf("instance: @G{%s}\n", stat.InstanceID)
		fmt.Printf("binding:  @G{%s}\n", stat.BindingID)
		fmt.Printf("status:   @C{%s}\n", stat.Status)
		os.Exit(0)

	case "unbind":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s} [@W{options}] @M{BINDING}\n\n", os.Args[0], command)
			fmt.Printf("  -s, --service  The name or ID of the service that the instance\n")
			fmt.Printf("                 was provisioned from.  This is required if the\n")
			fmt.Printf("                 instance details are not found in ~/.osbrc.\n")
			fmt.Printf("\n")
			fmt.Printf("  -p, --plan     The name or ID of the plan that the instance\n")
			fmt.Printf("                 was provisioned from.  This is required if the\n")
			fmt.Printf("                 instance details are not found in ~/.osbrc.\n")
			fmt.Printf("\n")
			fmt.Printf("  -i, --id       The ID of the service instance this binding belongs\n")
			fmt.Printf("                 to.  This is required if the binding details are not\n")
			fmt.Printf("                 found in ~/.osbrc.\n")
			fmt.Printf("\n")
			os.Exit(0)
		}

		unbinding(args)

		instance, service, plan, err := store.GetBindingDetails(c.URL, args[0])
		bail(err)

		if instance == "" {
			instance = opt.Unbind.ID
			if instance == "" {
				fmt.Fprintf(os.Stderr, "@R{instance not found in local ~/.osbrc}\n")
				fmt.Fprintf(os.Stderr, "You must specify the --id flag to the unbind operation.\n")
				os.Exit(1)
			}
		}
		if service == "" || plan == "" {
			catalog, err := c.GetCatalog()
			bail(err)

			if service == "" {
				service = opt.Unbind.Service
				if service == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", instance)
					fmt.Fprintf(os.Stderr, "You must specify the --service flag to the unbind operation.\n")
					os.Exit(1)
				}
			}
			if plan == "" {
				plan = opt.Unbind.Plan
				if plan == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", instance)
					fmt.Fprintf(os.Stderr, "You must specify the --plan flag to the unbind operation.\n")
					os.Exit(1)
				}
			}

			s, p, err := catalog.FindPlan(service, plan)
			bail(err)
			service = s
			plan = p
		}

		stat, err := c.Unbind(api.UnbindSpec{
			InstanceID: instance,
			BindingID:  args[0],
			ServiceID:  service,
			PlanID:     plan,
		})
		bail(err)

		store.RemoveBinding(c.URL, stat.InstanceID, stat.BindingID)
		if err := store.Write(opt.Data); err != nil {
			fmt.Fprintf(os.Stderr, "@Y{!!! %s}\n", err)
		}

		if opt.JSON {
			jsonify(stat)
			os.Exit(0)
		}

		fmt.Printf("instance: @G{%s}\n", stat.InstanceID)
		fmt.Printf("binding:  @G{%s}\n", stat.BindingID)
		fmt.Printf("status:   @C{%s}\n", stat.Status)
		os.Exit(0)

	case "deprovision":
		if opt.Help {
			fmt.Printf("USAGE: @G{%s} [@W{options}] @C{%s}\n\n", os.Args[0], command)
			os.Exit(0)
		}

		deprovisioning(args)

		instance := args[0]
		service, plan, _ := store.GetInstanceDetails(c.URL, instance)
		if service == "" || plan == "" {
			catalog, err := c.GetCatalog()
			bail(err)

			if service == "" {
				service = opt.Deprovision.Service
				if service == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", instance)
					fmt.Fprintf(os.Stderr, "You must specify the --service flag to the deprovision operation.\n")
					os.Exit(1)
				}
			}
			if plan == "" {
				plan = opt.Deprovision.Plan
				if plan == "" {
					fmt.Fprintf(os.Stderr, "@R{instance '%s' not found in local ~/.osbrc}\n", instance)
					fmt.Fprintf(os.Stderr, "You must specify the --plan flag to the deprovision operation.\n")
					os.Exit(1)
				}
			}

			s, p, err := catalog.FindPlan(service, plan)
			bail(err)
			service = s
			plan = p
		}
		stat, err := c.Deprovision(api.DeprovisionSpec{
			InstanceID: args[0],
			ServiceID:  service,
			PlanID:     plan,
		})
		bail(err)

		store.RemoveInstance(c.URL, args[0])
		if err := store.Write(opt.Data); err != nil {
			fmt.Fprintf(os.Stderr, "@Y{!!! %s}\n", err)
		}

		fmt.Printf("instance: @G{%s}\n", args[0])
		fmt.Printf("status:   @M{%s}\n", stat.Status)
		if stat.Operation != "" {
			fmt.Printf("operation: @C{%s}\n", stat.Operation)
		}
		os.Exit(0)
	}
}

func jsonify(x interface{}) {
	b, err := json.Marshal(x)
	bail(err)
	fmt.Printf("%s\n", string(b))
}

func connecting() {
	if opt.Endpoint == "" {
		fmt.Fprintf(os.Stderr, "@Y{missing required --endpoint flag or $OSB_URL environment variable}\n")
		os.Exit(1)
	}
	if opt.Username == "" {
		fmt.Fprintf(os.Stderr, "@Y{missing required --username flag or $OSB_USERNAME environment variable}\n")
		os.Exit(1)
	}
	if opt.Password == "" {
		fmt.Fprintf(os.Stderr, "@Y{missing required --password flag or $OSB_PASSWORD environment variable}\n")
		os.Exit(1)
	}
}

func provisioning(args []string) {
	connecting()
	if len(args) != 1 {
		fmt.Printf("USAGE: @Y{%s} [@W{options}] @C{provision} [-i INSTANCE-ID] SERVICE/PLAN\n", os.Args[0])
		os.Exit(1)
	}
}

func binding(args []string) {
	connecting()
	if len(args) != 1 {
		fmt.Printf("USAGE: @Y{%s} [@W{options}] @C{bind} [-i BINDING-ID] INSTANCE-ID\n", os.Args[0])
		os.Exit(1)
	}
}

func unbinding(args []string) {
	connecting()
	if len(args) != 1 {
		fmt.Printf("USAGE: @Y{%s} [@W{options}] @C{unbind} BINDING-ID\n", os.Args[0])
		os.Exit(1)
	}
}

func deprovisioning(args []string) {
	connecting()
	if len(args) != 1 {
		fmt.Printf("USAGE: @Y{%s} [@W{options}] @C{deprovision} INSTANCE-ID\n", os.Args[0])
		os.Exit(1)
	}
}
