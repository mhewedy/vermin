wget -qO- https://binaries.cockroachdb.com/cockroach-v19.2.6.linux-amd64.tgz | tar  xvz
sudo cp -i cockroach-v19.2.6.linux-amd64/cockroach /usr/local/bin/


for node in 0 1 2
do
	cockroach start \
	--insecure \
	--store=node_$node \
	--listen-addr=localhost:$((26257+node)) \
	--http-addr=localhost:$((8080+node)) \
	--join=localhost:26257,localhost:26258,localhost:26259 \
  --background > ~/cockroach.log 2>&1
done

cockroach init --insecure
