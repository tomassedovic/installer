data "openstack_images_image_v2" "bootstrap_image" {
  name        = "${var.image_name}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = "${var.flavor_name}"
}

data "ignition_systemd_unit" "haproxy_unit" {
    name = "bootkube-haproxy.service"
    enabled = true
    content = "[Service]\nType=oneshot\nExecStartPre=/sbin/setenforce 0\nExecStart=/bin/podman run -ti -d --name haproxy --net=host -v /etc/haproxy:/usr/local/etc/haproxy:ro docker.io/library/haproxy:1.7\n\n[Install]\nWantedBy=multi-user.target"
}

data "ignition_file" "haproxy_conf" {
    filesystem = "root"
    path = "/etc/haproxy/haproxy.cfg"
    source {
      source = "data:,listen%20ostest-api-80%0D%0A%20%20%20%20bind%200.0.0.0%3A80%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20ostest-bootstrap%20ostest-bootstrap.shiftstack.com%3A80%20check%0D%0A%20%20%20%20server%20ostest-master-0%20ostest-master-0.shiftstack.com%3A80%20check%0D%0A%20%20%20%20server%20ostest-master-1%20ostest-master-1.shiftstack.com%3A80%20check%0D%0A%20%20%20%20server%20ostest-master-2%20ostest-master-2.shiftstack.com%3A80%20check%0D%0A%0D%0Alisten%20ostest-api-6443%0D%0A%20%20%20%20bind%200.0.0.0%3A6443%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20ostest-bootstrap%20ostest-bootstrap.shiftstack.com%3A6443%20check%0D%0A%20%20%20%20server%20ostest-master-0%20ostest-master-0.shiftstack.com%3A6443%20check%0D%0A%20%20%20%20server%20ostest-master-1%20ostest-master-1.shiftstack.com%3A6443%20check%0D%0A%20%20%20%20server%20ostest-master-2%20ostest-master-2.shiftstack.com%3A6443%20check%0D%0A%0D%0Alisten%20ostest-api-443%0D%0A%20%20%20%20bind%200.0.0.0%3A443%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20ostest-bootstrap%20ostest-bootstrap.shiftstack.com%3A443%20check%0D%0A%20%20%20%20server%20ostest-master-0%20ostest-master-0.shiftstack.com%3A443%20check%0D%0A%20%20%20%20server%20ostest-master-1%20ostest-master-1.shiftstack.com%3A443%20check%0D%0A%20%20%20%20server%20ostest-master-2%20ostest-master-2.shiftstack.com%3A443%20check%0D%0A%0D%0Alisten%20ostest-api-49500%0D%0A%20%20%20%20bind%200.0.0.0%3A49500%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20ostest-bootstrap%20ostest-bootstrap.shiftstack.com%3A49500%20check%0D%0A%20%20%20%20server%20ostest-master-0%20ostest-master-0.shiftstack.com%3A49500%20check%0D%0A%20%20%20%20server%20ostest-master-1%20ostest-master-1.shiftstack.com%3A49500%20check%0D%0A%20%20%20%20server%20ostest-master-2%20ostest-master-2.shiftstack.com%3A49500%20check"
    }
} 

data "ignition_user" "core" {
    name = "core"
    ssh_authorized_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDX0EAM8x9uTjYWD+yZolNDuFxDbmy1cpmDOecL7+SxwoI84LuAGQnwxFofpnmKpNa6XQlYi1OSY2NQmMhrp8dl4e7+utc7ShqjFvTXBsHtQAOsboAWq7vL6fgwwEADUiWi+aERhtNJHjOz1EPOyN40M9yEui9I3SQnOQBFPmjMhFpy561m2qDs8LyoB3XMsmkRKmrLsTmYWWtf3abMsVfPjsZfN87oKJBQbYrfvOXoQ3wOa/IXvmCB2rf360LlHh0WiV1xFLggFdj659/huoPGs2B58op7Cep1YHprBvZTivetnGYhQWbha4WUh9UzLJtvxdG5mHzPRZcg71yeH8dv root@localhost.localdomain"
    ]
}

data "ignition_config" "config" {
  files = [
    "${data.ignition_file.haproxy_conf.id}"
  ]

  systemd = [
      "${data.ignition_systemd_unit.haproxy_unit.id}",
  ]
 
  users = [
      "${data.ignition_user.core.id}",
  ]
}

resource "openstack_objectstorage_object_v1" "lb_ignition" {
  container_name = "${var.swift_container}"
  name           = "load-balancer.ign"
  content        = "${data.ignition_config.config.rendered}" 
}

resource "openstack_objectstorage_tempurl_v1" "lb_ignition_tmpurl" {
  container = "${var.swift_container}"
  method    = "get"
  object    = "${openstack_objectstorage_object_v1.lb_ignition.name}"
  ttl       = 3600
}

data "ignition_config" "lb_redirect" {
  replace {
    source = "${openstack_objectstorage_tempurl_v1.lb_ignition_tmpurl.url}"
  }
}

resource "openstack_compute_instance_v2" "load_balancer" {
  name      = "${var.cluster_name}-api"
  flavor_id = "${data.openstack_compute_flavor_v2.bootstrap_flavor.id}"
  image_id  = "${data.openstack_images_image_v2.bootstrap_image.id}"

  user_data = "${data.ignition_config.lb_redirect.rendered}"

  network {
    port = "${var.lb_port_id}"
  }

  metadata {
    Name = "${var.cluster_name}-bootstrap"

    # "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    tectonicClusterID = "${var.cluster_id}"
  }
}
