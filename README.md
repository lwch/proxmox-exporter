# proxmox-exporter

proxmox exporter of prometheus

## node metrics

    # HELP node_cpu_frequency node cpu frequency of each core
    # TYPE node_cpu_frequency gauge
    node_cpu_frequency{node_name="pve",processor="0"} 1600
    node_cpu_frequency{node_name="pve",processor="1"} 1600
    # HELP node_cpu_loadavg node cpu load average
    # TYPE node_cpu_loadavg gauge
    node_cpu_loadavg{minute="1",node_name="pve"} 0.45
    node_cpu_loadavg{minute="15",node_name="pve"} 1.02
    node_cpu_loadavg{minute="5",node_name="pve"} 0.92
    # HELP node_cpu_usage node cpu usage ratio(precent)
    # TYPE node_cpu_usage gauge
    node_cpu_usage{node_name="pve"} 2.69030695310328
    # HELP node_info node info, labels:\nmodel: cpu model\nsockets: cpu sockets count\ncores: cpu cores\nthreads: cpu threads\nmhz: cpu frequency\nkernel_version: linux kernel version\npve_version: proxmox version
    # TYPE node_info gauge
    node_info{cores="6",kernel_version="Linux 5.15.39-3-pve #2 SMP PVE 5.15.39-3 (Wed, 27 Jul 2022 13:45:39 +0200)",mhz="1600.000",model="Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz",node_name="pve",pve_version="pve-manager/7.2-7/d0dd0e85",sockets="1",threads="12"} 10206
    # HELP node_iowait node iowait ratio(precent)
    # TYPE node_iowait gauge
    node_iowait{node_name="pve"} 0.20251757137751702
    # HELP node_memory_free free memory bytes of this node
    # TYPE node_memory_free gauge
    node_memory_free{node_name="pve"} 2.006562816e+10
    # HELP node_memory_total total memory bytes of this node
    # TYPE node_memory_total gauge
    node_memory_total{node_name="pve"} 3.3363566592e+10
    # HELP node_memory_used used memory bytes of this node
    # TYPE node_memory_used gauge
    node_memory_used{node_name="pve"} 1.3297938432e+10
    # HELP node_netin node received bytes
    # TYPE node_netin gauge
    node_netin{node_name="pve"} 8044.83333333333
    # HELP node_netout node sent bytes
    # TYPE node_netout gauge
    node_netout{node_name="pve"} 14875.6233333333
    # HELP node_rootfs_free free rootfs bytes of this node
    # TYPE node_rootfs_free gauge
    node_rootfs_free{node_name="pve"} 1.80893761536e+11
    # HELP node_rootfs_total total rootfs bytes of this node
    # TYPE node_rootfs_total gauge
    node_rootfs_total{node_name="pve"} 2.14698033152e+11
    # HELP node_rootfs_used used rootfs bytes of this node
    # TYPE node_rootfs_used gauge
    node_rootfs_used{node_name="pve"} 3.3804271616e+10
    # HELP node_sensors use sensors command to get device temperature and cpu fan speed
    # TYPE node_sensors gauge
    node_sensors{chip_name="acpitz-acpi-0",feature_name="temp1_input",label_name="temp1",node_name="pve"} -263.2
    node_sensors{chip_name="acpitz-acpi-0",feature_name="temp2_crit",label_name="temp2",node_name="pve"} 119
    node_sensors{chip_name="acpitz-acpi-0",feature_name="temp2_input",label_name="temp2",node_name="pve"} 77
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp1_crit",label_name="Package id 0",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp1_crit_alarm",label_name="Package id 0",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp1_input",label_name="Package id 0",node_name="pve"} 78
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp1_max",label_name="Package id 0",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp2_crit",label_name="Core 0",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp2_crit_alarm",label_name="Core 0",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp2_input",label_name="Core 0",node_name="pve"} 76
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp2_max",label_name="Core 0",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp3_crit",label_name="Core 1",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp3_crit_alarm",label_name="Core 1",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp3_input",label_name="Core 1",node_name="pve"} 74
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp3_max",label_name="Core 1",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp4_crit",label_name="Core 2",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp4_crit_alarm",label_name="Core 2",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp4_input",label_name="Core 2",node_name="pve"} 78
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp4_max",label_name="Core 2",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp5_crit",label_name="Core 3",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp5_crit_alarm",label_name="Core 3",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp5_input",label_name="Core 3",node_name="pve"} 76
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp5_max",label_name="Core 3",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp6_crit",label_name="Core 4",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp6_crit_alarm",label_name="Core 4",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp6_input",label_name="Core 4",node_name="pve"} 76
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp6_max",label_name="Core 4",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp7_crit",label_name="Core 5",node_name="pve"} 100
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp7_crit_alarm",label_name="Core 5",node_name="pve"} 0
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp7_input",label_name="Core 5",node_name="pve"} 76
    node_sensors{chip_name="coretemp-isa-0000",feature_name="temp7_max",label_name="Core 5",node_name="pve"} 100
    node_sensors{chip_name="nvme-pci-3a00",feature_name="temp1_alarm",label_name="Composite",node_name="pve"} 0
    node_sensors{chip_name="nvme-pci-3a00",feature_name="temp1_crit",label_name="Composite",node_name="pve"} 84.85
    node_sensors{chip_name="nvme-pci-3a00",feature_name="temp1_input",label_name="Composite",node_name="pve"} 56.85
    node_sensors{chip_name="nvme-pci-3a00",feature_name="temp1_max",label_name="Composite",node_name="pve"} 82.85
    node_sensors{chip_name="nvme-pci-3a00",feature_name="temp1_min",label_name="Composite",node_name="pve"} -273.15
    node_sensors{chip_name="pch_cannonlake-virtual-0",feature_name="temp1_input",label_name="temp1",node_name="pve"} 72
    # HELP node_smart_power_cycles power cycles of smart data
    # TYPE node_smart_power_cycles gauge
    node_smart_power_cycles{device="nvme0n1",node_name="pve",type="nvme"} 63
    node_smart_power_cycles{device="sda",node_name="pve",type="sata"} 3
    # HELP node_smart_poweron_hours poweron hours of smart data
    # TYPE node_smart_poweron_hours gauge
    node_smart_poweron_hours{device="nvme0n1",node_name="pve",type="nvme"} 10743
    node_smart_poweron_hours{device="sda",node_name="pve",type="sata"} 1
    # HELP node_smart_readden readden bytes of smart data(lba 512 bytes padding)
    # TYPE node_smart_readden gauge
    node_smart_readden{device="nvme0n1",node_name="pve",type="nvme"} 2.1106878e+07
    # HELP node_smart_temperature temperature of smart data
    # TYPE node_smart_temperature gauge
    node_smart_temperature{device="nvme0n1",node_name="pve",type="nvme"} 56.89999999999998
    node_smart_temperature{device="sda",node_name="pve",type="sata"} 40
    # HELP node_smart_used_percent used percent of smart data(nvme)
    # TYPE node_smart_used_percent gauge
    node_smart_used_percent{device="nvme0n1",node_name="pve",type="nvme"} 2
    # HELP node_smart_written written bytes of smart data(lba 512 bytes padding)
    # TYPE node_smart_written gauge
    node_smart_written{device="nvme0n1",node_name="pve",type="nvme"} 2.5995899e+07
    node_smart_written{device="sda",node_name="pve",type="sata"} 1.378059346e+09
    # HELP node_storage_free node storage free bytes
    # TYPE node_storage_free gauge
    node_storage_free{node_name="pve",storage_name="backup"} 4.97467867136e+11
    node_storage_free{node_name="pve",storage_name="local"} 1.80893761536e+11
    node_storage_free{node_name="pve",storage_name="local-lvm"} 1.65635413771e+11
    # HELP node_storage_info node storage info, labels:\ncontent_*: allowed content type\nstorage: storage name\ntype: storage type
    # TYPE node_storage_info gauge
    node_storage_info{content_backup="false",content_images="false",content_iso="true",content_rootdir="false",content_snippets="false",content_vztmpl="true",node_name="pve",storage="local",type="dir"} 10438
    node_storage_info{content_backup="false",content_images="true",content_iso="false",content_rootdir="true",content_snippets="false",content_vztmpl="false",node_name="pve",storage="local-lvm",type="lvmthin"} 10438
    node_storage_info{content_backup="true",content_images="false",content_iso="false",content_rootdir="false",content_snippets="false",content_vztmpl="false",node_name="pve",storage="backup",type="dir"} 10438
    # HELP node_storage_total node storage total bytes
    # TYPE node_storage_total gauge
    node_storage_total{node_name="pve",storage_name="backup"} 5.36608768e+11
    node_storage_total{node_name="pve",storage_name="local"} 2.14698033152e+11
    node_storage_total{node_name="pve",storage_name="local-lvm"} 6.442450944e+11
    # HELP node_storage_usage node storage usage ratio(precent)
    # TYPE node_storage_usage gauge
    node_storage_usage{node_name="pve",storage_name="backup"} 7.2941224963361
    node_storage_usage{node_name="pve",storage_name="local"} 15.7450308788193
    node_storage_usage{node_name="pve",storage_name="local-lvm"} 74.289999999882
    # HELP node_storage_used node storage used bytes
    # TYPE node_storage_used gauge
    node_storage_used{node_name="pve",storage_name="backup"} 3.9140900864e+10
    node_storage_used{node_name="pve",storage_name="local"} 3.3804271616e+10
    node_storage_used{node_name="pve",storage_name="local-lvm"} 4.78609680629e+11
    # HELP node_swap_free free swap bytes of this node
    # TYPE node_swap_free gauge
    node_swap_free{node_name="pve"} 8.57970688e+09
    # HELP node_swap_total total swap bytes of this node
    # TYPE node_swap_total gauge
    node_swap_total{node_name="pve"} 8.589930496e+09
    # HELP node_swap_used used swap bytes of this node
    # TYPE node_swap_used gauge
    node_swap_used{node_name="pve"} 1.0223616e+07
    # HELP node_uptime node uptime
    # TYPE node_uptime gauge
    node_uptime{node_name="pve"} 322277

## vm metrics

    # HELP vm_cpu_total vm max cpu core count
    # TYPE vm_cpu_total gauge
    vm_cpu_total{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 2
    # HELP vm_cpu_usage vm cpu usage ratio(precent)
    # TYPE vm_cpu_usage gauge
    vm_cpu_usage{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 0.132567768344343
    # HELP vm_disk_read vm disk readen bytes
    # TYPE vm_disk_read gauge
    vm_disk_read{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 7.43112704e+08
    # HELP vm_disk_total vm disk total bytes
    # TYPE vm_disk_total gauge
    vm_disk_total{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 2.0957446144e+10
    # HELP vm_disk_used vm disk used bytes
    # TYPE vm_disk_used gauge
    vm_disk_used{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 5.229576192e+09
    # HELP vm_disk_write vm disk written bytes
    # TYPE vm_disk_write gauge
    vm_disk_write{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 2.287190016e+09
    # HELP vm_info vm info, labels:\nuptime: uptime\ncore: cpu cores\nmemory: max memory bytes\ndisk: max disk bytes
    # TYPE vm_info gauge
    vm_info{core="1",disk="53687091200",memory="536870912",node_name="pve",uptime="0",vm_id="lxc/108",vm_name="backup",vm_status="stopped",vm_type="lxc"} 1
    # HELP vm_memory_total vm memory total bytes
    # TYPE vm_memory_total gauge
    vm_memory_total{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 1.073741824e+09
    # HELP vm_memory_used vm memory used bytes
    # TYPE vm_memory_used gauge
    vm_memory_used{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 9.250816e+07
    # HELP vm_netin vm network received bytes
    # TYPE vm_netin gauge
    vm_netin{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 1.026262351e+09
    # HELP vm_netout vm network sent bytes
    # TYPE vm_netout gauge
    vm_netout{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 9.54059148e+08
    # HELP vm_uptime vm uptime
    # TYPE vm_uptime gauge
    vm_uptime{node_name="pve",vm_id="lxc/100",vm_name="gw",vm_status="running",vm_type="lxc"} 322260