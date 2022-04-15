# Systemd 配置、安装和启动

- [Systemd 配置、安装和启动](#systemd-配置安装和启动)
	- [1. 前置操作（需要 root 权限）](#前置操作需要-root-权限)
	- [2. 创建 bifrost systemd unit 模板文件](#创建-bifrost-systemd-unit-模板文件)
	- [3. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)](#复制-systemd-unit-模板文件到-sysmted-配置目录需要有root权限)
	- [4. 启动 systemd 服务](#启动-systemd-服务)

## 1. 前置操作（需要 root 权限）

1. 根据注释配置 `environment.sh`

2. 创建 data 目录 

```
mkdir -p ${BIFROST_DATA_DIR}/bifrost
```

3. 创建 bin 目录，并将 `bifrost` 可执行文件复制过去

```bash
source ./environment.sh
mkdir -p ${BIFROST_INSTALL_DIR}/bin
cp bifrost ${BIFROST_INSTALL_DIR}/bin
```

4. 将 `bifrost` 配置文件拷贝到 `${BIFROST_CONFIG_DIR}` 目录下

```bash
mkdir -p ${BIFROST_CONFIG_DIR}
cp bifrost.yaml ${BIFROST_CONFIG_DIR}
```

## 2. 创建 bifrost systemd unit 模板文件

执行如下 shell 脚本生成 `bifrost.service.template`

```bash
source ./environment.sh
cat > bifrost.service.template <<EOF
[Unit]
Description=BIFROST
Documentation=https://github.com/ClessLi/bifrost/blob/master/init/README.md

[Service]
WorkingDirectory=${BIFROST_DATA_DIR}/bifrost
ExecStart=${BIFROST_INSTALL_DIR}/bin/bifrost --config=${BIFROST_CONFIG_DIR}/bifrost.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 3. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)

```bash
cp bifrost.service.template /etc/systemd/system/bifrost.service
```

## 4. 启动 systemd 服务

```bash
systemctl daemon-reload && systemctl enable bifrost && systemctl restart bifrost
```
