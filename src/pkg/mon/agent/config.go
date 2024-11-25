//nolint:lll,tagalign // Tags of config params are manually aligned.
package agent

import (
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

	"cspage/pkg/config"
	"cspage/pkg/http"
	"cspage/pkg/msg"
)

type vm struct {
	ID   string `json:"id"   param:"cloud-vm-id"   default:"" desc:"cloud instance ID"`
	Name string `json:"name" param:"cloud-vm-name" default:"" desc:"cloud instance name"`
	Type string `json:"type" param:"cloud-vm-type" default:"" desc:"cloud instance type"`
}

type env struct {
	config.BaseEnv

	AgentID string `json:"agent_id"`
	Cloud   string `json:"cloud"`
	Region  string `json:"cloud_region"  param:"cloud-region"  default:"" desc:"cloud region"`
	Zone    string `json:"cloud_zone"    param:"cloud-zone"    default:"" desc:"cloud availability zone"`
	VM      vm     `json:"cloud_vm"      param:"cloud-vm-*"`
}

type PingConfig struct {
	Count          int64         `param:"ping-count"           default:"10"`
	Interval       time.Duration `param:"ping-interval"        default:"500ms"`
	Timeout        time.Duration `param:"ping-timeout"         default:"10s"`
	PongInterval   time.Duration `param:"ping-pong-interval"   default:"250ms"`
	PongTimeout    time.Duration `param:"ping-pong-timeout"    default:"5s"`
	ResolveDelay   time.Duration `param:"ping-resolve-delay"   default:"10s"`
	ResolveTimeout time.Duration `param:"ping-resolve-timeout" default:"5s"`
}

type Config struct {
	config.BaseConfig         `param:"config-*"`
	http.HTTPConfig           `param:"http-*"`
	msg.PubsubPublisherConfig `param:"pubsub-*"`

	Env                      env           `param:"env-*"`
	PingInterval             time.Duration `param:"ping-interval"               default:"15s"   desc:"how often to send ping messages"`
	ProbeIntervalDefault     time.Duration `param:"probe-interval-default"      default:"120s"  desc:"default frequency for running monitoring probes"`
	ProbeLongIntervalDefault time.Duration `param:"probe-long-interval-default" default:"1200s" desc:"default frequency for running expensive monitoring probes"`
	ProbeTimeout             time.Duration `param:"probe-timeout"               default:"60s"   desc:"timeout for probe operations (cloud API calls)"`
	ProbeLongTimeout         time.Duration `param:"probe-long-timeout"          default:"300s"  desc:"timeout for long-running probe operations (cloud API calls)"`
	PubsubPingTopic          string        `param:"pubsub-ping-topic"  default:"mon-pings"   desc:"Pub/Sub topic for PING and AGENT messages"`
	PubsubProbeTopic         string        `param:"pubsub-probe-topic" default:"mon-probes"  desc:"Pub/Sub topic for PROBE messages"`
	SiteURL                  string        `param:"site-url"    default:"https://cloudstatus.page/api" desc:"API URL"`
	SiteSecret               string        `param:"site-secret" default:"Y2xvdWRzdGF0dXM6bmJ1c3IxMjM=" desc:"HTTP basic auth"`

	VPCIntraPing PingConfig `param:"common-vpc-intra-*" prefix:"common-vpc-intra-"`
	VPCInterPing PingConfig `param:"common-vpc-inter-*" prefix:"common-vpc-inter-"`
}

type Dummy struct {
	AgentPingProbeInterval    time.Duration `param:"dummy-agent-ping-probe-interval"    default:"-1s"`
	InternetPingProbeInterval time.Duration `param:"dummy-internet-ping-probe-interval" default:"-1s"`
}

type GCP struct {
	ProjectID string `param:"gcp-project-id" default:"cloudstatus-probe-t"`

	ComputeVMType                    string        `param:"gcp-compute-vm-type"           default:"e2-micro"`
	ComputeVMZonesSkip               []string      `param:"gcp-compute-vm-zone-skip"      default:""`
	ComputeVMDiskImage               string        `param:"gcp-compute-vm-disk-image"     default:"projects/debian-cloud/global/images/family/debian-12"`
	ComputeVMSubnetwork              string        `param:"gcp-compute-vm-subnetwork"     default:"projects/cloudstatus-probe-t/regions/europe-west3/subnetworks/europe-west3"`
	ComputeVMPrefix                  string        `param:"gcp-compute-vm-prefix"         default:"test-europe-west3"`
	ComputeVMProbeInterval           time.Duration `param:"gcp-compute-vm-probe-interval" default:"-1s"`
	ComputeVMSpotType                string        `param:"gcp-compute-vm-spot-type"           default:"e2-micro"`
	ComputeVMSpotZonesSkip           []string      `param:"gcp-compute-vm-spot-zone-skip"      default:""`
	ComputeVMSpotPrefix              string        `param:"gcp-compute-vm-spot-prefix"         default:"test-spot-europe-west3"`
	ComputeVMSpotProbeInterval       time.Duration `param:"gcp-compute-vm-spot-probe-interval" default:"-1s"`
	ComputeVMMetadataProbeInterval   time.Duration `param:"gcp-compute-vm-metadata-probe-interval" default:"-1s"`
	ComputeDiskSnapshotDiskName      string        `param:"gcp-compute-disk-snapshot-disk-name"      default:""`
	ComputeDiskSnapshotPrefix        string        `param:"gcp-compute-disk-snapshot-prefix"         default:"test-europe-west3"`
	ComputeDiskSnapshotProbeInterval time.Duration `param:"gcp-compute-disk-snapshot-probe-interval" default:"-1s"`

	StorageObjectBucketName    string        `param:"gcp-storage-object-bucket-name"    default:"cloudstatus-probe-t-europe-west3"`
	StorageObjectPrefix        string        `param:"gcp-storage-object-prefix"         default:"test/"`
	StorageObjectProbeInterval time.Duration `param:"gcp-storage-object-probe-interval" default:"-1s"`
	StorageBucketPrefix        string        `param:"gcp-storage-bucket-prefix"         default:"cloudstatus-probe-t-europe-west3-test"`
	StorageBucketProbeInterval time.Duration `param:"gcp-storage-bucket-probe-interval" default:"-1s"`

	PubsubProject         string        `param:"gcp-pubsub-project"          default:"cloudstatus-probe-t"`
	PubsubTopic           string        `param:"gcp-pubsub-topic"            default:"mon-probe-europe-west3"`
	PubsubSubscription    string        `param:"gcp-pubsub-subscription"     default:"mon-probe-europe-west3"`
	PubsubMessageInterval time.Duration `param:"gcp-pubsub-message-interval" default:"-1s"`

	VPCInterPingInterval time.Duration `param:"gcp-vpc-inter-ping-interval" default:"-1s"`
}

type AWS struct {
	EC2VMType                   string        `param:"aws-ec2-vm-type"             default:"t3.nano"`
	EC2VMZonesSkip              []string      `param:"aws-ec2-vm-zone-skip"        default:""`
	EC2VMDiskImageOwner         string        `param:"aws-ec2-vm-disk-image-owner" default:"amazon"`
	EC2VMDiskImageName          string        `param:"aws-ec2-vm-disk-image-name"  default:"al2023-ami-minimal-*-x86_64"`
	EC2VMVPCID                  string        `param:"aws-ec2-vm-vpc-id"           default:""`
	EC2VMPrefix                 string        `param:"aws-ec2-vm-prefix"           default:"test"`
	EC2VMProbeInterval          time.Duration `param:"aws-ec2-vm-probe-interval"   default:"-1s"`
	EC2VMSpotZonesSkip          []string      `param:"aws-ec2-vm-spot-zone-skip"   default:""`
	EC2VMSpotType               string        `param:"aws-ec2-vm-spot-type"             default:"t3.nano"`
	EC2VMSpotDiskImageName      string        `param:"aws-ec2-vm-spot-disk-image-name"  default:"al2023-ami-minimal-*-x86_64"`
	EC2VMSpotPrefix             string        `param:"aws-ec2-vm-spot-prefix"           default:"test-spot"`
	EC2VMSpotProbeInterval      time.Duration `param:"aws-ec2-vm-spot-probe-interval"   default:"-1s"`
	EC2VMMetadataProbeInterval  time.Duration `param:"aws-ec2-vm-metadata-probe-interval" default:"-1s"`
	EC2EBSSnapshotVolumeID      string        `param:"aws-ec2-ebs-snapshot-volume-id"      default:""`
	EC2EBSSnapshotPrefix        string        `param:"aws-ec2-ebs-snapshot-prefix"         default:"test"`
	EC2EBSSnapshotProbeInterval time.Duration `param:"aws-ec2-ebs-snapshot-probe-interval" default:"-1s"`

	S3ObjectBucketName    string        `param:"aws-s3-object-bucket-name"    default:"cloudstatus-probe-t-eu-central-1"`
	S3ObjectPrefix        string        `param:"aws-s3-object-prefix"         default:"test/"`
	S3ObjectProbeInterval time.Duration `param:"aws-s3-object-probe-interval" default:"-1s"`
	S3BucketPrefix        string        `param:"aws-s3-bucket-prefix"         default:"cloudstatus-probe-t-eu-central-1-test"`
	S3BucketProbeInterval time.Duration `param:"aws-s3-bucket-probe-interval" default:"-1s"`

	SQSQueueName       string        `param:"aws-sqs-queue-name"       default:"mon-probe"`
	SQSMessageInterval time.Duration `param:"aws-sqs-message-interval" default:"-1s"`

	VPCInterPingInterval time.Duration `param:"aws-vpc-inter-ping-interval" default:"-1s"`
}

type Azure struct {
	SubscriptionID string `param:"azure-subscription-id" default:""`
	ResourceGroup  string `param:"azure-resource-group"  default:"mon-probe-t-westeurope"`

	ComputeVMType                   string        `param:"azure-compute-vm-type"                  default:"Standard_B1ls"`
	ComputeVMDiskImagePublisher     string        `param:"azure-compute-vm-disk-image-publisher"  default:"canonical"`
	ComputeVMDiskImageOffer         string        `param:"azure-compute-vm-disk-image-offer"      default:"ubuntu-24_04-lts"`
	ComputeVMDiskImageSKU           string        `param:"azure-compute-vm-disk-image-sku"        default:"minimal"`
	ComputeVMDiskImageVersion       string        `param:"azure-compute-vm-disk-image-version"    default:"latest"`
	ComputeVMNICID                  string        `param:"azure-compute-vm-nic-id"                default:"/subscriptions/{uuid}/resourceGroups/mon-probe-t-westeurope/providers/Microsoft.Network/networkInterfaces/test-nic"`
	ComputeVMPrefix                 string        `param:"azure-compute-vm-prefix"                default:"test"`
	ComputeVMProbeInterval          time.Duration `param:"azure-compute-vm-probe-interval"        default:"-1s"`
	ComputeVMSpotType               string        `param:"azure-compute-vm-spot-type"              default:"Standard_B2ts_v2"`
	ComputeVMSpotNICID              string        `param:"azure-compute-vm-spot-nic-id"            default:"/subscriptions/{uuid}/resourceGroups/mon-probe-t-westeurope/providers/Microsoft.Network/networkInterfaces/test-spot-nic"`
	ComputeVMSpotPrefix             string        `param:"azure-compute-vm-spot-prefix"            default:"test-spot"`
	ComputeVMSpotProbeInterval      time.Duration `param:"azure-compute-vm-spot-probe-interval"    default:"-1s"`
	ComputeVMMetadataProbeInterval  time.Duration `param:"azure-compute-vm-metadata-probe-interval" default:"-1s"`
	ComputeVHDSnapshotDiskID        string        `param:"azure-compute-vhd-snapshot-disk-id"        default:""`
	ComputeVHDSnapshotPrefix        string        `param:"azure-compute-vhd-snapshot-prefix"         default:"test"`
	ComputeVHDSnapshotProbeInterval time.Duration `param:"azure-compute-vhd-snapshot-probe-interval" default:"-1s"`

	StorageAccountName            string        `param:"azure-storage-account-name"        default:"probetwesteurope"`
	StorageBlobContainerName      string        `param:"azure-storage-blob-container-name" default:"objects"`
	StorageBlobPrefix             string        `param:"azure-storage-blob-prefix"         default:"test/"`
	StorageBlobProbeInterval      time.Duration `param:"azure-storage-blob-probe-interval" default:"-1s"`
	StorageContainerPrefix        string        `param:"azure-storage-container-prefix"    default:"mon-probe-t-test"`
	StorageContainerProbeInterval time.Duration `param:"azure-storage-container-probe-interval" default:"-1s"`

	ServiceBusNamespace            string        `param:"azure-servicebus-namespace"              default:"probetwesteurope"`
	ServiceBusQueueName            string        `param:"azure-servicebus-queue-name"             default:"mon-probe"`
	ServiceBusQueueMessageInterval time.Duration `param:"azure-servicebus-queue-message-interval" default:"-1s"`

	VPCInterPingInterval time.Duration `param:"azure-vpc-inter-ping-interval" default:"-1s"`
}

type Cloud interface {
	Dummy | GCP | AWS | Azure
}

type CloudConfig[T Cloud] struct {
	Config `param:"config-*"`

	Cloud T `param:"cloud-*"`
}

type (
	GCPConfig   = CloudConfig[GCP]
	AWSConfig   = CloudConfig[AWS]
	AzureConfig = CloudConfig[Azure]
)

//nolint:gochecknoglobals // Common function.
var Die = config.Die

//nolint:gochecknoglobals // Common function.
var DieLog = config.DieLog

func NewConfig[T Cloud]() *CloudConfig[T] {
	// Using UUID version 6 to keep node ID in there; This way we can detect when an agent was
	// restarted on an existing node or when an agent was started on a completely new node.
	agentID, err := uuid.NewV6()
	if err != nil {
		Die("Could not generate agent ID", "err", err)
	}

	cfg := &CloudConfig[T]{
		Config: Config{
			Env: env{
				BaseEnv: config.NewBaseEnv(),
				AgentID: agentID.String(),
				Cloud:   strings.ToLower(reflect.TypeOf(*new(T)).Name()),
			},
		},
	}
	config.InitConfig(cfg, &cfg.BaseConfig)

	return cfg
}

func (c *Config) ProbeInterval(d time.Duration) time.Duration {
	if d < 0 {
		return c.ProbeIntervalDefault
	}
	return d
}

func (c *Config) ProbeLongInterval(d time.Duration) time.Duration {
	if d < 0 {
		return c.ProbeLongIntervalDefault
	}
	return d
}
