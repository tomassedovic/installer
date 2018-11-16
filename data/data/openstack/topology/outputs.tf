output "bootstrap_port_id" {
  value = "${openstack_networking_port_v2.bootstrap_port.id}"
}

output "bootstrap_port_fixed_ip" {
  value = "${openstack_networking_port_v2.bootstrap_port.all_fixed_ips[0]}"
}

output "master_sg_id" {
  value = "${openstack_networking_secgroup_v2.master.id}"
}

output "master_subnet_ids" {
  value = "${local.master_subnet_ids}"
}
