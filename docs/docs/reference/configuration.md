# Configuration

## Config

Package `config` contains all the options available for configuring a Talos cluster.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`clusterName`</td>
<td markdown="1">string</td>
<td markdown="1">Configures the cluster's name.<details><summary>*Show example*</summary>
```yaml
clusterName: my-cluster
```
</details></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`endpoint`</td>
<td markdown="1">string</td>
<td markdown="1"><details><summary>Configures the cluster's controlplane endpoint.</summary>Can be an IP address or a DNS hostname</details><details><summary>*Show example*</summary>
```yaml
endpoint: https://192.168.200.10:6443
```
</details></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`nodes`</td>
<td markdown="1">[][Node](#node)</td>
<td markdown="1">List of nodes configurations<details><summary>*Show example*</summary>
```yaml
nodes:
  - hostname: kmaster1
    ipAddress: 192.168.200.11
    controlPlane: true
    installDiskSelector:
      size: 128GB
  - hostname: kworker1
    ipAddress: 192.168.200.12
    controlPlane: false
    installDisk: /dev/sda
    networkInterfaces:
      - interface: eth0
        dhcp: true
```
</details></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`talosVersion`</td>
<td markdown="1">string</td>
<td markdown="1"><details><summary>Talos version to perform the installation.</summary>Image reference for each Talos release can be found on <br />[Talos GitHub release page](https://github.com)</details><details><summary>*Show example*</summary>
```yaml
talosVersion: v1.5.2
```
</details></td>
<td markdown="1" align="center">`"latest"`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`kubernetesVersion`</td>
<td markdown="1">string</td>
<td markdown="1">Allows for supplying the Kubernetes version to use.</details><details><summary>*Show example*</summary>
```yaml
kubernetesVersion: v1.28.1
```
</details></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`domain`</td>
<td markdown="1">string</td>
<td markdown="1">Allows for supplying the domain used by Kubernetes DNS.</details><details><summary>*Show example*</summary>
```yaml
domain: mycluster.com
```
</details></td>
<td markdown="1" align="center">`"cluster.local"`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`allowSchedulingOnMasters`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to allow running workload on controlplane nodes.</details><details><summary>*Show example*</summary>
```yaml
allowSchedulingOnMasters: true
```
</details></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`allowSchedulingOnControlPlanes`</td>
<td markdown="1">bool</td>
<td markdown="1"><details><summary>Whether to allow running workload on controlplane nodes.</summary>It is an alias to `allowSchedulingOnMasters`</details><details><summary>*Show example*</summary>
```yaml
allowSchedulingOnControlPlanes: true
```
</details></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`additionalMachineCertSans`</td>
<td markdown="1">[]string</td>
<td markdown="1">Extra certificate SANs for the machine's certificate.<details><summary>*Show example*</summary>
```yaml
additionalMachineCertSans:
  - 10.0.0.10
  - 172.16.0.10
  - 192.168.0.10
```
</details></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`additionalApiServerCertSans`</td>
<td markdown="1">[]string</td>
<td markdown="1">Extra certificate SANs for the API server's certificate.<details><summary>*Show example*</summary>
```yaml
additionalApiServerCertSans:
  - 1.2.3.4
  - 4.5.6.7
  - mycluster.local
```
</details></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`cniConfig`</td>
<td markdown="1">[CNIConfig](#cniconfig)</td>
<td markdown="1">The CNI to be used for the cluster's network.<details><summary>*Show example*</summary>
```yaml
cniConfig:
  name: custom
  urls:
    - https://docs.projectcalico.org/archive/v3.20/manifests/canal.yaml
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`imageFactory`</td>
<td markdown="1">[ImageFactory](#imagefactory)</td>
<td markdown="1">Configures selfhosted image factory.<details><summary>*Show example*</summary>
```yaml
imageFactory:
  registryURL: myfactory.com
  schematicEndpoint: /schematics
  protocol: https
  installerURLTmpl: {{.RegistryURL}}/installer/{{.ID}}:{{.Version}}
  ISOURLTmpl: {{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}/{{.Arch}}.iso
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`patches`</td>
<td markdown="1">[]string</td>
<td markdown="1"><details><summary>Patches to be applied to all nodes.</summary>List of strings containing RFC6902 JSON patches, strategic merge patches,<br />or a file containing them</details><details><summary>*Show example*</summary>
```yaml
patches:
  - |-
    - op: add
      path: /machine/kubelet/extraArgs
      value:
        rotate-server-certificates: "true"
  - |-
    machine:
      env:
        MYENV: value
  - "@./a-patch.yaml"
```
</details></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`controlPlane`</td>
<td markdown="1">[NodeConfigs](#nodeconfigs)</td>
<td markdown="1">Configurations targetted for all controlplane nodes.</details><details><summary>*Show example*</summary>
```yaml
controlPlane:
  kernelModules:
    - name: br_netfilter
      parameters:
        - nf_conntrack_max=131072
  patches:
    - |-
      - op: add
        path: /machine/kubelet/extraArgs
        value:
          rotate-server-certificates: "true"
    - |-
      machine:
        env:
          MYENV: value
    - "@./a-patch.yaml"
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`worker`</td>
<td markdown="1">[NodeConfigs](#nodeconfigs)</td>
<td markdown="1">Configurations targetted for all worker nodes.</details><details><summary>*Show example*</summary>
```yaml
worker:
  kernelModules:
    - name: br_netfilter
      parameters:
        - nf_conntrack_max=131072
  patches:
    - |-
      - op: add
        path: /machine/kubelet/extraArgs
        value:
          rotate-server-certificates: "true"
    - |-
      machine:
        env:
          MYENV: value
    - "@./a-patch.yaml"
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## Node

`Node` defines machine configurations for each node.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`hostname`</td>
<td markdown="1">string</td>
<td markdown="1">Configures the hostname of a node.<details><summary>*Show example*</summary>
```yaml
hostname: kmaster1
```
</details></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`ipAddress`</td>
<td markdown="1">string</td>
<td markdown="1"><details><summary>IP address where the node can be reached, can be IP or comma separated list of IPs.</summary>Needed for endpoint and node address inside `talosconfig`.</details><details><summary>*Show example*</summary>
```yaml
ipAddress: 192.168.200.11
```
</summary></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`installDisk`</td>
<td markdown="1">string</td>
<td markdown="1">The disk used for installation.<details><summary>*Show example*</summary>
```yaml
installDisk: /dev/sda
```
</summary></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`installDiskSelector`</td>
<td markdown="1">[InstallDiskSelector](#installdiskselector)</td>
<td markdown="1"><details><summary>Look up disk used for installation.</summary>Required if `installDisk` is not specified.</details><details><summary>*Show example*</summary>
```yaml
installDiskSelector:
  size: 128GB
  model: WDC*
  name: /sys/block/sda/device/name
  busPath: /pci0000:00/0000:00:17.0/ata1/host0/target0:0:0/0:0:0:0
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`controlPlane`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether the node is a controlplane.<details><summary>*Show example*</summary>
```yaml
controlPlane: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`ignoreHostname`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to set `machine.network.hostname` to the generated config file.<details><summary>*Show example*</summary>
```yaml
ignoreHostname: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`overridePatches`</td>
<td markdown="1">bool</td>
<td markdown="1"><details><summary>Whether `patches` defined here should override the one defined in node group.</summary>By default they will get appended instead.</details><details><summary>*Show example*</summary>
```yaml
overridePatches: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`overrideExtraManifests`</td>
<td markdown="1">bool</td>
<td markdown="1"><details><summary>Whether `extraManifests` defined here should override the one defined in node group.</summary>By default they will get appended instead.</details><details><summary>*Show example*</summary>
```yaml
overrideExtraManifests: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">-</td>
<td markdown="1">[NodeConfigs](#nodeconfigs)</td>
<td markdown="1">Node specific configurations that will override node group configurations.<details><summary>*Show example*</summary>
```yaml
talosImageURL: factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e
nodeLabels:
  rack: rack1a
nodeTaints:
  exampleTaint: exampletaintValue:NoSchedule
disableSearchDomain: true
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## NodeConfigs

`NodeConfigs` defines machine configurations.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`talosImageURL`</td>
<td markdown="1">string</td>
<td markdown="1">Allows for supplying the node level image used to perform the installation.<details><summary>*Show example*</summary>
```yaml
talosImageURL: factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e
```
</details></td>
<td markdown="1" align="center">`""`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`machineSpec`</td>
<td markdown="1">[MachineSpec](#machinespec)</td>
<td markdown="1"><details><summary>Machine hardware specification for the node.</summary>Only used for `genurl iso` subcommand.</details><details><summary>*Show example*</summary>
```yaml
machineSpec:
  mode: metal
  arch: arm64
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`ingressFirewall`</td>
<td markdown="1">[IngressFirewall](#ingressfirewall)</td>
<td markdown="1">Machine firewall specification for the node.<details><summary>*Show example*</summary>
```yaml
ingressFirewall:
  defaultAction: block
  rules:
    - name: kubelet-ingress
      portSelector:
        ports:
          - 10250
        protocol: tcp
      ingress:
        - subnet: 172.20.0.0/24
          except: 172.20.0.1/32
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`extensionServices`</td>
<td markdown="1">[][ExtensionService](#extensionservice)</td>
<td markdown="1">Machine extension services specification for the node.<details><summary>*Show example*</summary>
```yaml
extensionServices:
  - name: nut-client
    configFiles:
      - content: MONITOR upsmonHost 1 remote pass password
        mountPath: /usr/local/etc/nut/upsmon.conf
    environment:
      - UPS_NAME=ups
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`nodeLabels`</td>
<td markdown="1">map[string]string</td>
<td markdown="1">Labels to be added to the node.<details><summary>*Show example*</summary>
```yaml
rack: rack1a
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`nodeAnnotations`</td>
<td markdown="1">map[string]string</td>
<td markdown="1">Annotations to be added to the node.<details><summary>*Show example*</summary>
```yaml
rack: rack1a
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`nodeTaints`</td>
<td markdown="1">map[string]string</td>
<td markdown="1">Node taints for the node.<details><summary>*Show example*</summary>
```yaml
exampleTaint: exampleTaintValue:NoSchedule
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`disableSearchDomain`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to disable generating default search domain.<details><summary>*Show example*</summary>
```yaml
disableSearchDomain: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`machineDisks`</td>
<td markdown="1">[][MachineDisk](#machinedisk)</td>
<td markdown="1">List of additional disks to partition, format, mount.<details><summary>*Show example*</summary>
```yaml
machineDisks:
  - device: /dev/disk/by-id/ata-CT500MX500SSD1_2149E5EC1D9D
    partitions:
      - mountpoint: /var/mnt/sata
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`noSchematicValidate`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to skip schematic validation.<details><summary>*Show example*</summary>
```yaml
noSchematicValidate: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`disableSearchDomain`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to disable generating default search domain.<details><summary>*Show example*</summary>
```yaml
disableSearchDomain: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`machineFiles`</td>
<td markdown="1">[][MachineFile](#machinefile)</td>
<td markdown="1">List of additional files to create inside the node.<details><summary>*Show example*</summary>
```yaml
machineFiles:
  - content: |
      TS_AUTHKEY=123456
    permissions: 0o644
    path: /var/etc/tailscale/auth.env
    op: create
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`schematic`</td>
<td markdown="1">[Schematic](#schematic)</td>
<td markdown="1">Configure Talos image customization to be used in the installer image<details><summary>*Show example*</summary>
```yaml
schematic:
  customization:
    extraKernelArgs:
      - net.ifnames=0
    systemExtensions:
      officialExtensions:
        - siderolabs/intel-ucode
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`isoSchematic`</td>
<td markdown="1">[Schematic](#schematic)</td>
<td markdown="1">Configure Talos image customization to be used for ISO image<details><summary>*Show example*</summary>
```yaml
isoSchematic:
  customization:
    extraKernelArgs:
      - net.ifnames=0
    systemExtensions:
      officialExtensions:
        - siderolabs/intel-ucode
```
</summary></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`kernelModules`</td>
<td markdown="1">[][KernelModuleConfig](#kernelmoduleconfig)</td>
<td markdown="1">List of additional kernel modules to load.<details><summary>*Show example*</summary>
```yaml
kernelModules:
  - name: br_netfilter
    parameters:
      - nf_conntrack_max=131072
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`nameservers`</td>
<td markdown="1">[]string</td>
<td markdown="1">List of nameservers for the node.<details><summary>*Show example*</summary>
```yaml
nameservers:
  - 8.8.8.8
  - 1.1.1.1
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`networkInterfaces`</td>
<td markdown="1">[][Device](#device)</td>
<td markdown="1">List of network interface configurations for the node.<details><summary>*Show example*</summary>
```yaml
networkInterfaces:
  - interface: enp0s1
    addresses:
      - 192.168.2.0/24
    routes:
      - network: 0.0.0.0/0
        gateway: 192.168.2.1
        metric: 1024
    mtu: 1500
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`extraManifests`</td>
<td markdown="1">[]string</td>
<td markdown="1">List of manifest files to be added for the node.<details><summary>*Show example*</summary>
```yaml
extraManifests:
  - etcd-firewall.yaml
  - kubelet-firewall.yaml
```
</summary></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`patches`</td>
<td markdown="1">[]string</td>
<td markdown="1"><details><summary>Patches to be applied to the node.</summary>List of strings containing RFC6902 JSON patches, strategic merge patches,<br />or a file containing them.</details><details><summary>*Show example*</summary>
```yaml
patches:
  - |-
    - op: add
      path: /machine/kubelet/extraArgs
      value:
        rotate-server-certificates: "true"
  - |-
    machine:
      env:
        MYENV: value
  - "@./a-patch.yaml"
```
</details></td>
<td markdown="1" align="center">`[]`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## ImageFactory

`ImageFactory` defines configuration for selfhosted image-factory.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`registryURL`</td>
<td markdown="1">string</td>
<td markdown="1">Registry URL of the factory.<details><summary>*Show example*</summary>
```yaml
registryURL: myfactory.com
```
</details></td>
<td markdown="1" align="center">`"factory.talos.dev"`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`protocol`</td>
<td markdown="1">string</td>
<td markdown="1">Protocol the registry is listening to.<details><summary>*Show example*</summary>
```yaml
protocol: http
```
</summary></td>
<td markdown="1" align="center">`https`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`schematicEndpoint`</td>
<td markdown="1">string</td>
<td markdown="1">Path to do HTTP POST request to the registry.</details><details><summary>*Show example*</summary>
```yaml
schematicEndpoint: /schematics
```
</summary></td>
<td markdown="1" align="center">`/schematics`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`installerURLTmpl`</td>
<td markdown="1">string</td>
<td markdown="1"><details><summary>Go template to parse the full installer URL.</summary>Available placeholders: `RegistryURL`,`ID`,`Version`, `Secureboot`</details><details><summary>*Show example*</summary>
```yaml
installerURLTmpl: "{{.RegistryURL}}/installer/{{.ID}}:{{.Version}}"
```
</summary></td>
<td markdown="1" align="center">`{{.RegistryURL}}/installer{{if .Secureboot}}-secureboot{{end}}/{{.ID}}:{{.Version}}`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`ISOURLTmpl`</td>
<td markdown="1">string</td>
<td markdown="1"><details><summary>Go template to parse the full ISO image URL.</summary>Available placeholders: `Protocol`,`RegistryURL`,`ID`,`Version`,`Mode`,`Arch`, `Secureboot`, `UseUKI`</details><details><summary>*Show example*</summary>
```yaml
ISOURLTmpl: "{{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}-{{.Arch}}.iso"
```
</summary></td>
<td markdown="1" align="center">`{{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}-{{.Arch}}{{if .Secureboot}}-secureboot{{end}}{{if and .Secureboot .UseUKI}}-uki.efi{{else}}.iso{{end}}`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## MachineSpec

`MachineSpec` defines machine hardware configurations for a node.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`mode`</td>
<td markdown="1">string</td>
<td markdown="1">Machine mode.<details><summary>*Show example*</summary>
```yaml
mode: metal
```
</details></td>
<td markdown="1" align="center">`"metal"`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`arch`</td>
<td markdown="1">string</td>
<td markdown="1">Machine architecture.<details><summary>*Show example*</summary>
```yaml
arch: arm64
```
</summary></td>
<td markdown="1" align="center">`amd64`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`secureboot`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to enable Secure Boot.<details><summary>*Show example*</summary>
```yaml
secureboot: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`useUKI`</td>
<td markdown="1">bool</td>
<td markdown="1">Whether to use UKI if Secure Boot is enabled.<details><summary>*Show example*</summary>
```yaml
useUKI: true
```
</summary></td>
<td markdown="1" align="center">`false`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## IngressFirewall

`IngressFirewall` defines machine firewall configuration for a node.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`defaultAction`</td>
<td markdown="1">`string`</td>
<td markdown="1"><details><summary>Default action for all not explicitly configured traffic.</summary>Can be "accept" or "block"</details><details><summary>*Show example*</summary>
```yaml
defaultAction: accept
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`rules`</td>
<td markdown="1">[][NetworkRule](#networkrule)</td>
<td markdown="1"><details><summary>List of matching network rules to allow or block against the defaultAction.</summary>If `defaultAction` is set to block, matching rules will be allowed vice versa.</details><details><summary>*Show example*</summary>
```yaml
rules:
  - name: kubelet-ingress
    portSelector:
      ports:
        - 10250
      protocol: tcp
    ingress:
      - subnet: 172.20.0.0/24
        except: 172.20.0.1/32
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

</table>

## ExtensionService

`ExtensionService` defines machine extension service configuration for a node.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`name`</td>
<td markdown="1">`string`</td>
<td markdown="1">Name of the extension service config.<details><summary>*Show example*</summary>
```yaml
name: nut-client
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`configFiles`</td>
<td markdown="1">[][ConfigFile](#configfile)</td>
<td markdown="1">The config files for the extension service.<details><summary>*Show example*</summary>
```yaml
configFiles:
  - content: MONITOR upsmonHost 1 remote pass password
    mountPath: /usr/local/etc/nut/upsmon.conf
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`environment`</td>
<td markdown="1">[]string</td>
<td markdown="1">The environment for the extension service.<details><summary>*Show example*</summary>
```yaml
environment:
  - UPS_NAME=ups
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:negative_squared_cross_mark:</td>
</tr>

</table>

## NetworkRule

`NetworkRule` defines the firewall rules to match.

<table markdown="1">
<tr markdown="1">
<th markdown="1">Field</th><th>Type</th><th>Description</th><th>Default Value</th><th>Required</th>
</tr>

<tr markdown="1">
<td markdown="1">`name`</td>
<td markdown="1">`string`</td>
<td markdown="1">Name of the rule.<details><summary>*Show example*</summary>
```yaml
name: kubelet-ingress
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`portSelector`</td>
<td markdown="1">[PortSelector](#portselector)</td>
<td markdown="1">Ports and protocols on the host affected by the rule.<details><summary>*Show example*</summary>
```yaml
portSelector:
  ports:
    - 10250
  protocol: tcp
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

<tr markdown="1">
<td markdown="1">`ingress`</td>
<td markdown="1">[][IngressConfig](#ingressconfig)</td>
<td markdown="1">List of source subnets allowed to access the host ports/protocols.<details><summary>*Show example*</summary>
```yaml
ingress:
  - subnet: 172.20.0.0/24
    except: 172.20.0.1/32
```
</details></td>
<td markdown="1" align="center">`nil`</td>
<td markdown="1" align="center">:white_check_mark:</td>
</tr>

</table>

## CNIConfig

`CNIConfig` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#cniconfig" target="_blank">`v1alpha1.CNIConfig`</a>

## InstallDiskSelector

`InstallDiskSelector` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#installdiskselector" target="_blank">`v1alpha1.InstallDiskSelector`</a>.

## MachineDisk

`MachineDisk` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#machinedisk" target="_blank">`v1alpha1.MachineDisk`</a>

## MachineFile

`MachineFile` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#machinefile" target="_blank">`v1alpha1.MachineFile`</a>

## InstallExtensionConfig

`InstallExtensionConfig` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#installextensionconfig" target="_blank">`v1alpha1.InstallExtensionConfig`</a>

## Schematic

`Schematic` is type of upstream Talos Image Factory <a href="https://pkg.go.dev/github.com/siderolabs/image-factory/pkg/schematic#Schematic" target="_blank">`schematic.Schematic`</a>

## KernelModuleConfig

`KernelModuleConfig` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#kernelmoduleconfig" target="_blank">`v1alpha1.KernelModuleConfig`</a>

## Device

`Device` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/#device" target="_blank">`v1alpha1.Device`</a>

## PortSelector

`PortSelector` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/network/networkruleconfig/#NetworkRuleConfig.portSelector" target="_blank">`network.RulePortSelector`</a>

## IngressConfig

`IngressConfig` is type of upstream Talos <a href="https://www.talos.dev/latest/reference/configuration/network/networkruleconfig/#NetworkRuleConfig.ingress" target="_blank">`network.IngressConfig`</a>

## ConfigFile

`ConfigFile` is type of upstream Talos <a href="https://www.talos.dev/v1.7/reference/configuration/extensions/extensionserviceconfig/#ExtensionServiceConfig.configFiles" target="_blank">`extensions.ConfigFile`</a>
