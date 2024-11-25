output "vpc_id" {
  value = module.vpc.vpc_id
}

output "vpc_subnet_ids" {
  value = module.vpc.public_subnets
}

output "security_group_id" {
  value = aws_security_group.mon-agent.id
}
