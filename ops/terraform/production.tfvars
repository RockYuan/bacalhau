bacalhau_version       = "v0.1.28"
bacalhau_port          = "1235"
bacalhau_connect_node0 = "QmdZQ7ZbhnvWY1J12XYKGHApJ6aufKyLNSvf8jZBrBaAVL"
ipfs_version           = "v0.12.2"
gcp_project            = "bacalhau-production"
instance_count         = 3
region                 = "us-east4"
zone                   = "us-east4-c"
volume_size_gb         = 500
boot_disk_size_gb      = 500
machine_type           = "e2-standard-16"
protect_resources      = true
auto_subnets           = true
ingress_cidrs          = ["0.0.0.0/0"]
ssh_access_cidrs       = ["0.0.0.0/0"]
