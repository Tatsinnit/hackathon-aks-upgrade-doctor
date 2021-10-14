module github.com/Tatsinnit/hackathon-aks-upgrade-doctor

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v58.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.21
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/gosuri/uilive v0.0.4 // indirect
	github.com/gosuri/uiprogress v0.0.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.2.1
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
)
