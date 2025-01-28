Установка всего необходимого:
- [ ] Установка docker
```
sudo apt update
sudo apt install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
apt-cache policy docker-ce
sudo apt install docker-ce
sudo systemctl status docker
sudo usermod -aG docker ${USER}
su - ${USER}
sudo usermod -aG docker vboxuser
```
- [ ] Установка docker-compose
```
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
```
- [ ] Установка kubernetes
```
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt update
sudo apt install kubeadm kubelet kubectl
sudo apt-mark hold kubeadm kubelet kubectl
kubeadm version
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
sudo nano /etc/modules-load.d/containerd.conf
# Вставить, сохранить:
overlay
br_netfilter

sudo modprobe overlay
sudo modprobe br_netfilter
sudo nano /etc/sysctl.d/kubernetes.conf
# Вставить, сохранить:
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1

sudo sysctl --system
sudo hostnamectl set-hostname master.k8s
sudo nano /etc/hosts
# Вставить свои данные, сохранить:
192.168.1.1   master.k8s
192.168.1.2   node1.k8s
192.168.1.3   node2.k8s

# Далее идет только для мастера
---
sudo nano /etc/default/kubelet
# Заменить на:
KUBELET_EXTRA_ARGS="--cgroup-driver=cgroupfs"

sudo systemctl daemon-reload && sudo systemctl restart kubelet
sudo nano /etc/docker/daemon.json
# Вставить, сохранить:
{
      "exec-opts": ["native.cgroupdriver=systemd"],
      "log-driver": "json-file",
      "log-opts": {
      "max-size": "100m"
   },
       "storage-driver": "overlay2"
       }
       
sudo systemctl daemon-reload && sudo systemctl restart docker
sudo nano /lib/systemd/system/kubelet.service.d/10-kubeadm.conf
# Добавить:
Environment="KUBELET_EXTRA_ARGS=--fail-swap-on=false"

sudo systemctl daemon-reload && sudo systemctl restart kubelet
sudo nano /etc/containerd/config.toml - закомментировать disabled_plugins = ["cri"] - #disabled_plugins = ["cri"]
sudo systemctl restart containerd
---

# Запустить главный узел:
sudo kubeadm init --pod-network-cidr=192.168.0.0/16
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
export KUBECONFIG=$HOME/.kube/config

# Установка calico:
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.29.1/manifests/tigera-operator.yaml
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.29.1/manifests/custom-resources.yaml
kubectl taint nodes --all node-role.kubernetes.io/control-plane-

# Запуск на нодах:
sudo scp yuriy@192.168.1.1:~/.kube/config ~/.kube/config
sudo chown $( id -u):$( id -g) $HOME/.kube/config 
export KUBECONFIG=~/.kube/config

Заменить троеточие после sudo kubeadm ниже, на вывод команды - kubeadm token create --print-join-command
sudo kubeadm ...
```
- [ ] [Файл для создания, запуска всего необходимого](https://github.com/manzhelayv/myvmessenger/blob/main/kubernetes/job.sh), где "~/mnt" - смонтированная папка.
- [ ] После установки для работы коннектора кафки:
``` 
kubectl get po -n database
kubectl exec -it "name mongodb pod" mongosh -n database
rs.initiate({
    _id: "rs0",
     members: [
        { _id: 0, host: "mongodb.database:27017" }
     ]
})
curl --request POST \
  --url http://localhost:30077/connectors \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.3.0' \
  --data '{
    "name": "mongo-connector",
    "config": {
        "connector.class" : "io.debezium.connector.mongodb.MongoDbConnector",
        "mongodb.connection.string" : "mongodb://mongodb.database:27017/?replicaSet=rs0",
        "topic.prefix" : "chat",
			  "collection.include.list" : "test.chat",
			  "signal.kafka.groupId": "debezium-chat-group-id"
    }
}'
```
- [ ] Все вместе например с Ubuntu Desktop на VirtualBox нужно минимум 80 гигабайт на диске и 16 гигабайт памяти.
- [ ] [Backend](https://github.com/manzhelayv/myvmessenger/tree/main/backend). Go get клиентов с приватного gitlab, нужно [заменить на](https://github.com/manzhelayv/myvmessenger/tree/main/backend/client)
- [ ] [Frontend](https://github.com/manzhelayv/myvmessenger/tree/main/frontend)
- [ ] [Для кубернетеса](https://github.com/manzhelayv/myvmessenger/tree/main/kubernetes). Гитлаб раннер запускается через docker-compose, после создания runner-a в gitlab, нужно поменять данные в .env файле и перезапустить register-runner. Так же можно [развернуть гитлаб локально](https://github.com/manzhelayv/myvmessenger/tree/main/docker/gitlab)
- [ ] Развернуть все необходимое для бэкенда, кроме golang, можно с помощью [docker-compose](https://github.com/manzhelayv/myvmessenger/tree/main/docker)