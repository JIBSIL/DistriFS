# recommend running this on a machine with good bandwidth. You can get a free Google Cloud Shell instance for free!

echo "Welcome to the DistriFS benchmarking tool. Installation of an IPFS node, a DistriFS node and a torrenting server will be done in 5 seconds. CTRL+C to exit."
sleep 5

rm -rf /tmp/distrifs-benchmark

echo "Generating dummy data..."
mkdir /tmp/distrifs-benchmark && cd /tmp/distrifs-benchmark

fallocate -l 1G 1gb-file.bin
fallocate -l 100M 100mb-file.bin
fallocate -l 1M 1mb-file.bin

echo "Installing IPFS..."
wget https://dist.ipfs.tech/kubo/v0.28.0/kubo_v0.28.0_linux-amd64.tar.gz
tar -xvf kubo_v0.28.0_linux-amd64.tar.gz
cd kubo
sudo bash install.sh
ipfs init
cd /tmp/distrifs-benchmark
ipfs add 1mb-file.bin
ipfs add 100mb-file.bin
ipfs add 1gb-file.bin
ipfs daemon > /dev/null 2>&1 &

echo "Installing DistriFS..."
git clone https://github.com/JIBSIL/DistriFS
cd DistriFS
cd server
go install
go run . > /dev/null 2>&1 & 
cd ..
cd indexer
go install
go run . > /dev/null 2>&1 &

echo "Installing nginx..."
cd /tmp/distrifs-benchmark
sudo apt update
sudo apt-get install -y nginx
sudo mkdir /var/www/html/bench
sudo cp *.bin /var/www/html/bench
sudo service nginx start

echo "Installing torrent server..."
sudo apt-get install -y transmission-daemon
sudo cp *.bin /var/lib/transmission-daemon/downloads
sudo service transmission-daemon start

echo "If you are running this on a temporary server, run sudo ufw disable. If not, please do not disable your firewall and instead whitelist ports for nginx, DistriFS, IPFS and Transmission."