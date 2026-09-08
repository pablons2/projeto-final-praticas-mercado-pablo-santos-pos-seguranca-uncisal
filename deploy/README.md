# Deploy em produção — Oracle Cloud "Always Free" (passo a passo)

Alvo: IncidentTrack (Docker Compose) numa VM Oracle Cloud Always Free,
atendendo todos os itens do Eixo 1 do `Escopo_e_elementos_obrigatorios.md`.

Substitua em todo o guia:

| Marcador | Significado |
|---|---|
| `<IP>` | IP público reservado da VM |
| `<user>` | usuário SSH da imagem (`ubuntu`) |
| `<seu-email>` | e-mail para a conta Let's Encrypt |
| `DEPLOY_PATH` | `/home/ubuntu/incidenttrack` |

---

## 1. Criar a conta Oracle Cloud

1. <https://www.oracle.com/cloud/free/> → **Start for free**.
2. E-mail, país, telefone (SMS de verificação), **cartão de crédito**
   (verificação de identidade — não há cobrança no Always Free).
3. **Home Region:** escolha com cuidado — **não muda depois** e todos os
   recursos Always Free ficam nela. Prefira uma região próxima com
   capacidade (ex.: `sa-saopaulo-1` / `us-ashburn-1`).
4. Aguarde o provisionamento da tenancy (~5–15 min) e faça login no
   Console (<https://cloud.oracle.com>).

---

## 2. Criar a VM (Compute Instance)

**Menu ☰ → Compute → Instances → Create instance.**

1. **Name:** `incidenttrack`.
2. **Compartment:** o raiz (default) serve.
3. **Placement:** deixe o Availability Domain sugerido.
4. **Image and shape → Edit:**
   - **Image:** *Canonical Ubuntu* **24.04** (é a mais nova do catálogo
     Oracle; ver nota de PQC no passo 9).
   - **Shape:** `VM.Standard.A1.Flex` (ARM Ampere, Always Free) com
     **1 OCPU / 6 GB** — ou, se der "out of capacity", `VM.Standard.E2.1.Micro`
     (AMD, Always Free).
5. **Networking:**
   - **Create new VCN** (deixa o wizard criar VCN + subnet pública).
   - **Assign a public IPv4 address:** *Yes*.
6. **Add SSH keys:**
   - Cole sua **chave pública** (`cat ~/.ssh/id_ed25519.pub`) ou faça
     upload. Guarde a privada — é o único acesso.
7. **Create.** Quando o estado ficar *Running*, anote o **Public IP address**.

### 2.1. Reservar o IP (para não mudar)

**Instance → Resources → Attached VNICs → (a VNIC) → IPv4 Addresses →
edit do IP público → Reserved public IP → Create.** Assim o `<IP>` fica
fixo mesmo se a instância for reiniciada/parada.

---

## 3. Abrir as portas na nuvem (Security List)

Least privilege — só 22, 80, 443.

**Networking → Virtual Cloud Networks → (sua VCN) → Security Lists →
Default Security List → Add Ingress Rules:**

| Stateless | Source CIDR | IP Protocol | Dest. Port |
|---|---|---|---|
| No | `0.0.0.0/0` | TCP | 22 |
| No | `0.0.0.0/0` | TCP | 80 |
| No | `0.0.0.0/0` | TCP | 443 |

(Se existir uma regra ampla liberando tudo, remova.) Não é preciso mexer
em egress.

---

## 4. Primeiro acesso + hardening do SSH

```bash
ssh <user>@<IP>          # ex.: ssh ubuntu@<IP>

sudo apt update && sudo apt full-upgrade -y

# desabilita senha e login de root por SSH
sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/'  /etc/ssh/sshd_config
sudo sed -i 's/^#*PermitRootLogin.*/PermitRootLogin no/'                /etc/ssh/sshd_config
# a Oracle deixa um drop-in que reativa senha; neutralize:
sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config.d/*.conf 2>/dev/null || true
sudo systemctl restart ssh
```

Teste num **novo terminal** antes de fechar este: `ssh <user>@<IP>` deve
entrar pela chave; `ssh -o PubkeyAuthentication=no <user>@<IP>` deve ser
recusado.

---

## 5. Firewall do host (UFW)

A imagem Oracle usa `iptables` com regras próprias; o UFW por cima
simplifica e cumpre o requisito de firewall no host.

```bash
sudo apt install -y ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable
sudo ufw status verbose
```

> Se após habilitar o UFW o SSH cair, use o **Console Connection** da
> Oracle (Instance → Resources → Console connection) para recuperar.

---

## 6. Fail2Ban — SSH, 4 erros, ban de 24 h

```bash
sudo apt install -y fail2ban

sudo tee /etc/fail2ban/jail.d/sshd.local >/dev/null <<'EOF'
[sshd]
enabled  = true
port     = 22
maxretry = 4
findtime = 600
bantime  = 86400
EOF

sudo systemctl enable --now fail2ban
sudo fail2ban-client status sshd
```

`maxretry = 4` e `bantime = 86400` (24 h) são exatamente o que o escopo pede.

---

## 7. Docker + Compose v2

```bash
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker "$USER"
newgrp docker
docker compose version
```

---

## 8. Clonar o repositório

```bash
cd ~
git clone https://github.com/pablons2/projeto-final-praticas-mercado-pablo-santos-pos-seguranca-uncisal.git incidenttrack
cd incidenttrack        # este é o DEPLOY_PATH
```

---

## 9. TLS + PQC — a decisão de arquitetura

**Ubuntu 24.04 (imagem mais nova da Oracle) traz OpenSSL 3.0**, que
**não** tem o grupo pós-quântico `X25519MLKEM768`. O item "PQC ativado"
do Qualys exige OpenSSL **3.5+**. Duas saídas:

- **(A) Recomendada — terminar o TLS num contêiner nginx novo**
  (`nginx:1.29-alpine` = OpenSSL 3.5). Não mexe no SO, é reproduzível e
  já entra no `docker compose`. O "web server nginx" do requisito
  continua sendo nginx, só que containerizado.
- **(B) Atualizar o SO** para 25.10 (`sudo do-release-upgrade` — três
  saltos a partir de 24.04, ~1 h, com reboot) e usar nginx do host.

Este guia segue **(A)**.

### 9.1. Subir a stack da aplicação (sem porta pública ainda)

```bash
cp .env.production.example .env
nano .env
#   APP_ENV=production
#   COOKIE_SECURE=true
#   JWT_SECRET=<openssl rand -base64 48>     # gere na hora, exclusivo de produção
#   APP_PORT=8080
#   APP_ORIGIN=https://<IP>

docker compose up -d --build
curl -fsS http://127.0.0.1:8080/healthz      # -> ok
```

### 9.2. Emitir o certificado de IP (Certbot ≥ 5.4, webroot)

O plugin `--nginx` **não** suporta IP — usa-se `certonly --webroot`.

```bash
sudo snap install --classic certbot
sudo ln -sf /snap/bin/certbot /usr/bin/certbot
certbot --version        # precisa ser >= 5.4

sudo mkdir -p /var/www/certbot

# servidor temporário só para o desafio ACME na porta 80
sudo python3 -m http.server 80 --directory /var/www/certbot &
HTTP_PID=$!

sudo certbot certonly --webroot -w /var/www/certbot \
  --ip-address <IP> \
  --preferred-profile shortlived \
  --agree-tos -m <seu-email> --no-eff-email

sudo kill $HTTP_PID
```

Certificado em `/etc/letsencrypt/live/<IP>/`. É **short-lived (~6 dias)**.

### 9.3. Contêiner proxy TLS (nginx 3.5 + PQC)

O proxy **já está versionado** no repositório — nada de arquivo escrito à
mão na VM:

| Arquivo no repo | Papel |
|---|---|
| `deploy/proxy/templates/default.conf.template` | config nginx (redirect 80→443, TLS 1.2/1.3, `ssl_ecdh_curve X25519MLKEM768:…`, HSTS, OCSP, proxy → `web`) |
| `docker-compose.prod.yml` | serviço `proxy` (`nginx:1.29-alpine`), publica 80/443, remove a porta do `web`, monta `/etc/letsencrypt` e `/var/www/certbot` |

`${CERT_NAME}` no template é preenchido pelo entrypoint do nginx a partir
da variável `CERT_NAME` (host de `APP_ORIGIN`). O `deploy.yml` grava
`CERT_NAME=<IP>` no `.env` da VM automaticamente.

Subida manual (o pipeline faz o mesmo):

```bash
cd ~/incidenttrack
# .env já tem APP_ORIGIN=https://<IP>; garanta a linha CERT_NAME:
grep -q '^CERT_NAME=' .env || echo "CERT_NAME=<IP>" >> .env

docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
curl -fsS https://<IP>/healthz          # -> ok
curl -I  http://<IP>/                    # -> 301 para https
```

> A partir daqui, todo deploy (manual ou pelo pipeline) usa os **dois**
> arquivos. Enquanto `/etc/letsencrypt/live/<IP>/` não existir, o
> `deploy.yml` sobe só a stack base e avisa.

### 9.4. Renovação automática (cert de 6 dias → renovar a cada 6 h)

```bash
# recarrega o proxy sempre que o cert renova
sudo mkdir -p /etc/letsencrypt/renewal-hooks/deploy
sudo tee /etc/letsencrypt/renewal-hooks/deploy/reload-proxy.sh >/dev/null <<'EOF'
#!/bin/sh
cd /home/ubuntu/incidenttrack && \
  docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T proxy nginx -s reload
EOF
sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/reload-proxy.sh

# encurta o timer do snap para 6/6 h
sudo mkdir -p /etc/systemd/system/snap.certbot.renew.timer.d
sudo tee /etc/systemd/system/snap.certbot.renew.timer.d/override.conf >/dev/null <<'EOF'
[Timer]
OnCalendar=
OnCalendar=*-*-* 0/6:00:00
RandomizedDelaySec=30m
EOF
sudo systemctl daemon-reload
sudo systemctl restart snap.certbot.renew.timer

# teste (usa webroot; sobe um http.server temporário como no 9.2 se precisar)
sudo certbot renew --dry-run
```

---

## 10. GitHub — Environment + Secrets (liga o CD)

**Repo → Settings → Environments → New environment:** `production`.

**Repo → Settings → Secrets and variables → Actions → New repository secret:**

| Secret | Valor |
|---|---|
| `SSH_PRIVATE_KEY` | conteúdo de `~/.ssh/incidenttrack_deploy` (passo 11) |
| `SSH_KNOWN_HOSTS` | saída de `ssh-keyscan -H <IP>` |
| `SERVER_HOST` | `<IP>` |
| `SERVER_USER` | `ubuntu` |
| `DEPLOY_PATH` | `/home/ubuntu/incidenttrack` |
| `JWT_SECRET` | o mesmo valor do `.env` da VM |
| `APP_PORT` | `8080` |
| `APP_ORIGIN` | `https://<IP>` |

> O `deploy.yml` faz `git reset --hard` e regenera o `.env` (com
> `CERT_NAME`) a partir desses Secrets. `docker-compose.prod.yml` e
> `deploy/proxy/` **são versionados** — o pipeline os aplica sozinho
> assim que `/etc/letsencrypt/live/<IP>/` existe na VM (bootstrap do §9).
> O que **não** fica no repo é o `.env` real e os certificados.

---

## 11. Chave SSH dedicada ao deploy

Na **sua máquina**:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/incidenttrack_deploy -C gh-actions -N ""
ssh-copy-id -i ~/.ssh/incidenttrack_deploy.pub <user>@<IP>
ssh-keyscan -H <IP>                 # cole em SSH_KNOWN_HOSTS
cat ~/.ssh/incidenttrack_deploy    # cole em SSH_PRIVATE_KEY (chave inteira)
```

---

## 12. Disparar o pipeline

```bash
git commit --allow-empty -m "ci: primeiro deploy de producao"
git push origin main
```

**Actions:** `CI` (lint/test/build) → ao passar na `main` → `CD`
(`workflow_run`) → SSH na VM → `docker compose up -d --build` →
healthcheck em `https://<IP>/healthz`.

---

## 13. Qualys SSL Labs

1. <https://www.ssllabs.com/ssltest/> → informe `<IP>` → *Do not show
   the results on the boards* → **Submit**.
2. Metas:
   - **Overall Rating: A**.
   - Em *Configuration → Protocols*: só **TLS 1.2 e 1.3**.
   - Em *Configuration → Key Exchange* / *Named Groups*: aparece
     **X25519MLKEM768** (ou `MLKEM768`) → prova o PQC.
   - HSTS presente.
3. Se cair de A: revise cadeia do certificado completa (`fullchain.pem`),
   OCSP stapling, e remova qualquer cipher fraco.

---

## 14. Fechamento (Eixo 2)

- [ ] **Revogar** o PAT `ghp_xTTb…` (<https://github.com/settings/tokens>) e usar um novo.
- [ ] Ativar **2FA** na conta `pablons2`.
- [ ] Confirmar o repositório **público** (Settings → General → Danger Zone).
- [ ] No `README.md`: adicionar `https://<IP>` (app no ar) e o print da nota Qualys.

---

## Referência rápida de operação

```bash
cd ~/incidenttrack
docker compose ps
docker compose logs -f --tail=100
docker compose up -d --build          # aplicar mudança manual
sudo fail2ban-client status sshd      # bans ativos
sudo certbot certificates             # validade do cert
```
