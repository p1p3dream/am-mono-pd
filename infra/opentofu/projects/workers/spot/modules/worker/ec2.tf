# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance.

data "aws_ami" "debian" {
  most_recent = true
  owners      = ["136693071363"]

  filter {
    name   = "architecture"
    values = ["arm64"]
  }

  filter {
    name = "name"
    # Prevents selecting backports images.
    values = ["*debian-${var.config.vars.debian_version}-arm64*"]
  }
}

resource "aws_key_pair" "user" {
  key_name   = var.user_keypair_key_name
  public_key = var.user_keypair_public_key
}

resource "aws_instance" "main" {
  ami           = data.aws_ami.debian.id
  instance_type = var.config.vars.ec2_instance_type
  key_name      = aws_key_pair.user.key_name
  monitoring    = true

  # Override the default for the workers-vpc.
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.spot_profile.name
  subnet_id                   = var.subnet_id
  vpc_security_group_ids = concat(var.vpc_security_group_ids, [
    aws_security_group.ssh.id,
  ])

  root_block_device {
    volume_type           = var.config.vars.block_devices.root.volume_type
    volume_size           = var.config.vars.block_devices.root.volume_size
    iops                  = var.config.vars.block_devices.root.iops
    throughput            = var.config.vars.block_devices.root.throughput
    delete_on_termination = var.config.vars.block_devices.root.delete_on_termination
  }

  metadata_options {
    # https://aws.amazon.com/blogs/security/get-the-full-benefits-of-imdsv2-and-disable-imdsv1-across-your-aws-infrastructure/.
    http_tokens = "required"
  }

  dynamic "instance_market_options" {
    for_each = var.config.vars.spot_instance ? [1] : []
    content {
      market_type = "spot"
      spot_options {
        instance_interruption_behavior = "terminate"
        spot_instance_type             = "one-time"
      }
    }
  }

  user_data_replace_on_change = false
  user_data_base64            = var.user_data_base64
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume.
resource "aws_ebs_volume" "persistent" {
  # Do not use "availability_zone = aws_instance.main.availability_zone" because
  # it'll replace the volume every time the instance is recreated.
  availability_zone = var.config.vars.ebs_availability_zone

  size       = var.config.vars.block_devices.persistent.volume_size
  type       = var.config.vars.block_devices.persistent.volume_type
  iops       = var.config.vars.block_devices.persistent.iops
  throughput = var.config.vars.block_devices.persistent.throughput

  encrypted = true
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/volume_attachment.html.
resource "aws_volume_attachment" "persistent_attachment" {
  device_name = var.config.vars.block_devices.persistent.device_name
  volume_id   = aws_ebs_volume.persistent.id
  instance_id = aws_instance.main.id

  # Prevents the instance from being terminated before the volume is detached.
  stop_instance_before_detaching = true
}
