# How to configure this for a client

1. Install Docker
2. Configure prometheus.yml with the api-key in line 13 (You can generate your API keys on your account via Profile Settings > Api Keys > Generate New Key)
3. Configure prometheus.yml with the subaccount in line 14.
4. (Optional)  if you want to scrape multiple users/subaccounts, you need to create jobs for each of them in prometheus.yml.
5. Run ```docker-compose up -d```
6. Go to http://localhost:3000/?orgId=1 (user admin, pass admin)
7. Add datasource prometheus in http://localhost:3000/datasources (url http://host.docker.internal:9090)
8. Import dashboard in http://localhost:3000/dashboard/import (check dashboard.json file in root)
9. Check http://localhost:9090/targets to check if scrape happened and was successful.
10. You should be able to see the metrics now.