# How to run this locally with grafana and prometheus

1. Go to dev folder with terminal
2. (Optional) modify prometheus.yml if you want to scrape another user than selfmining.
3. (Optional) run ```docker-compose build --no-cache``` to clean the build cache
4. Run ```docker-compose up -d```
5. Go to http://localhost:3000/?orgId=1 (user admin, pass admin)
6. Add datasource prometheus (url http://localhost:9090 and access 'browser')
7. Import dashboard (check dashboard.json file in root)
8. Check http://localhost:9090/targets to check if scrape happened and was successful
9. You should be able to see the metrics now.

# How to run locally for debugging or check export variables
1. Run/Debug main.go
2. In a browser go to this page for example http://localhost:8080/miner?pool=luxor&apikey={apikey}&subaccount=selfmining