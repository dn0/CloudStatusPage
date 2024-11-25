package vpc

import (
	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/msg"
)

const (
	pingProbeName = "azure_vpc" // doc="Azure Virtual Network"
)

//nolint:mnd // Documentation is used by `make db/sql`.
func NewIntraPingProbeJob(
	cfg *agent.AzureConfig,
	topic *msg.PubsubTopic,
	host string,
) *VPC.PingProbeJob[agent.Azure] {
	return VPC.NewPingProbeJob(cfg, topic, VPC.NewPingProbe[agent.Azure](
		&cfg.Config.VPCIntraPing,
		pingProbeName,
		host,
		host,
		&VPC.PingAction{
			Id:   91, // name="vpc.intra.Ping64B" doc="64 bytes ICMP ping within a virtual network subnet" url=""
			Size: VPC.PingSize64B,
		},
		&VPC.PingAction{
			Id:   92, // name="vpc.intra.Ping1500B" doc="1500 bytes ICMP ping within a virtual network subnet" url=""
			Size: VPC.PingSize1500B,
		},
	))
}

//nolint:lll,mnd,funlen // Documentation is used by `make db/sql`.
func NewInterPingProbe[T agent.Azure](cfg *agent.AzureConfig) *VPC.AgentPingProbe[T] {
	return VPC.NewAgentPingProbe[T](
		&cfg.Config,
		pingProbeName,
		map[string]*VPC.PingAction{
			"eastus":             {Id: 9201, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eastus" doc="64 bytes ICMP ping sent to virtual network in eastus" url=""
			"eastus2":            {Id: 9202, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eastus2" doc="64 bytes ICMP ping sent to virtual network in eastus2" url=""
			"eastus3":            {Id: 9203, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eastus3" doc="64 bytes ICMP ping sent to virtual network in eastus3" url=""
			"westus":             {Id: 9204, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westus" doc="64 bytes ICMP ping sent to virtual network in westus" url=""
			"westus2":            {Id: 9205, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westus2" doc="64 bytes ICMP ping sent to virtual network in westus2" url=""
			"westus3":            {Id: 9206, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westus3" doc="64 bytes ICMP ping sent to virtual network in westus3" url=""
			"centralus":          {Id: 9207, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B centralus" doc="64 bytes ICMP ping sent to virtual network in centralus" url=""
			"westcentralus":      {Id: 9208, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westcentralus" doc="64 bytes ICMP ping sent to virtual network in westcentralus" url=""
			"southcentralus":     {Id: 9209, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southcentralus" doc="64 bytes ICMP ping sent to virtual network in southcentralus" url=""
			"northcentralus":     {Id: 9210, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B northcentralus" doc="64 bytes ICMP ping sent to virtual network in northcentralus" url=""
			"canadaeast":         {Id: 9211, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B canadaeast" doc="64 bytes ICMP ping sent to virtual network in canadaeast" url=""
			"canadacentral":      {Id: 9212, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B canadacentral" doc="64 bytes ICMP ping sent to virtual network in canadacentral" url=""
			"mexicocentral":      {Id: 9213, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B mexicocentral" doc="64 bytes ICMP ping sent to virtual network in mexicocentral" url=""
			"brazilsouth":        {Id: 9231, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B brazilsouth" doc="64 bytes ICMP ping sent to virtual network in brazilsouth" url=""
			"brazilsoutheast":    {Id: 9232, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B brazilsoutheast" doc="64 bytes ICMP ping sent to virtual network in brazilsoutheast" url=""
			"chilecentral":       {Id: 9233, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B chilecentral" doc="64 bytes ICMP ping sent to virtual network in chilecentral" url=""
			"westeurope":         {Id: 9241, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westeurope" doc="64 bytes ICMP ping sent to virtual network in westeurope" url=""
			"northeurope":        {Id: 9242, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B northeurope" doc="64 bytes ICMP ping sent to virtual network in northeurope" url=""
			"ukwest":             {Id: 9243, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ukwest" doc="64 bytes ICMP ping sent to virtual network in ukwest" url=""
			"uksouth":            {Id: 9244, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B uksouth" doc="64 bytes ICMP ping sent to virtual network in uksouth" url=""
			"switzerlandwest":    {Id: 9245, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B switzerlandwest" doc="64 bytes ICMP ping sent to virtual network in switzerlandwest" url=""
			"switzerlandnorth":   {Id: 9246, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B switzerlandnorth" doc="64 bytes ICMP ping sent to virtual network in switzerlandnorth" url=""
			"germanywestcentral": {Id: 9247, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B germanywestcentral" doc="64 bytes ICMP ping sent to virtual network in germanywestcentral" url=""
			"germanynorth":       {Id: 9248, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B germanynorth" doc="64 bytes ICMP ping sent to virtual network in germanynorth" url=""
			"francesouth":        {Id: 9249, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B francesouth" doc="64 bytes ICMP ping sent to virtual network in francesouth" url=""
			"francecentral":      {Id: 9250, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B francecentral" doc="64 bytes ICMP ping sent to virtual network in francecentral" url=""
			"swedencentral":      {Id: 9251, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B swedencentral" doc="64 bytes ICMP ping sent to virtual network in swedencentral" url=""
			"swedensouth":        {Id: 9252, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B swedensouth" doc="64 bytes ICMP ping sent to virtual network in swedensouth" url=""
			"norwaywest":         {Id: 9253, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B norwaywest" doc="64 bytes ICMP ping sent to virtual network in norwaywest" url=""
			"norwayeast":         {Id: 9254, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B norwayeast" doc="64 bytes ICMP ping sent to virtual network in norwayeast" url=""
			"spaincentral":       {Id: 9255, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B spaincentral" doc="64 bytes ICMP ping sent to virtual network in spaincentral" url=""
			"italynorth":         {Id: 9256, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B italynorth" doc="64 bytes ICMP ping sent to virtual network in italynorth" url=""
			"polandcentral":      {Id: 9257, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B polandcentral" doc="64 bytes ICMP ping sent to virtual network in polandcentral" url=""
			"greececentral":      {Id: 9258, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B greececentral" doc="64 bytes ICMP ping sent to virtual network in greececentral" url=""
			"austriaeast":        {Id: 9259, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B austriaeast" doc="64 bytes ICMP ping sent to virtual network in austriaeast" url=""
			"belgiumcentral":     {Id: 9260, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B belgiumcentral" doc="64 bytes ICMP ping sent to virtual network in belgiumcentral" url=""
			"denmarkeast":        {Id: 9261, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B denmarkeast" doc="64 bytes ICMP ping sent to virtual network in denmarkeast" url=""
			"finlandcentral":     {Id: 9262, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B finlandcentral" doc="64 bytes ICMP ping sent to virtual network in finlandcentral" url=""
			"australiaeast":      {Id: 9271, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australiaeast" doc="64 bytes ICMP ping sent to virtual network in australiaeast" url=""
			"australiacentral":   {Id: 9272, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australiacentral" doc="64 bytes ICMP ping sent to virtual network in australiacentral" url=""
			"australiacentral2":  {Id: 9273, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australiacentral2" doc="64 bytes ICMP ping sent to virtual network in australiacentral2" url=""
			"australiasoutheast": {Id: 9274, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australiasoutheast" doc="64 bytes ICMP ping sent to virtual network in australiasoutheast" url=""
			"newzealandnorth":    {Id: 9275, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B newzealandnorth" doc="64 bytes ICMP ping sent to virtual network in newzealandnorth" url=""
			"eastasia":           {Id: 9276, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eastasia" doc="64 bytes ICMP ping sent to virtual network in eastasia" url=""
			"southeastasia":      {Id: 9277, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southeastasia" doc="64 bytes ICMP ping sent to virtual network in southeastasia" url=""
			"malaysiawest":       {Id: 9278, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B malaysiawest" doc="64 bytes ICMP ping sent to virtual network in malaysiawest" url=""
			"indonesiacentral":   {Id: 9279, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B indonesiacentral" doc="64 bytes ICMP ping sent to virtual network in indonesiacentral" url=""
			"koreasouth":         {Id: 9280, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B koreasouth" doc="64 bytes ICMP ping sent to virtual network in koreasouth" url=""
			"koreacentral":       {Id: 9281, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B koreacentral" doc="64 bytes ICMP ping sent to virtual network in koreacentral" url=""
			"japanwest":          {Id: 9282, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B japanwest" doc="64 bytes ICMP ping sent to virtual network in japanwest" url=""
			"japaneast":          {Id: 9283, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B japaneast" doc="64 bytes ICMP ping sent to virtual network in japaneast" url=""
			"westindia":          {Id: 9284, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B westindia" doc="64 bytes ICMP ping sent to virtual network in westindia" url=""
			"southindia":         {Id: 9285, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southindia" doc="64 bytes ICMP ping sent to virtual network in southindia" url=""
			"centralindia":       {Id: 9286, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B centralindia" doc="64 bytes ICMP ping sent to virtual network in centralindia" url=""
			"indiasouthcentral":  {Id: 9287, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B indiasouthcentral" doc="64 bytes ICMP ping sent to virtual network in indiasouthcentral" url=""
			"taiwannorth":        {Id: 9288, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B taiwannorth" doc="64 bytes ICMP ping sent to virtual network in taiwannorth" url=""
			"uaenorth":           {Id: 9289, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B uaenorth" doc="64 bytes ICMP ping sent to virtual network in uaenorth" url=""
			"uaecentral":         {Id: 9290, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B uaecentral" doc="64 bytes ICMP ping sent to virtual network in uaecentral" url=""
			"qatarcentral":       {Id: 9291, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B qatarcentral" doc="64 bytes ICMP ping sent to virtual network in qatarcentral" url=""
			"saudiarabiaeast":    {Id: 9292, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B saudiarabiaeast" doc="64 bytes ICMP ping sent to virtual network in saudiarabiaeast" url=""
			"israelcentral":      {Id: 9293, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B israelcentral" doc="64 bytes ICMP ping sent to virtual network in israelcentral" url=""
			"southafricawest":    {Id: 9294, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southafricawest" doc="64 bytes ICMP ping sent to virtual network in southafricawest" url=""
			"southafricanorth":   {Id: 9295, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southafricanorth" doc="64 bytes ICMP ping sent to virtual network in southafricanorth" url=""
		},
	)
}
