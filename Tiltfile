# docker build 

docker_build("live-epg",'.',dockerfile="Dockerfile")

yaml = helm("./charts/live-epg","live-epg",set=[
    "mariadb.enabled=true",
    "image.repository=live-epg",
    "image.tag=latest",
    "app.env.data.DATASOURCE_MYSQL_DNS=live_epg:changeme@tcp(live-epg-mariadb:3306)/live_epg?charset=utf8mb4&parseTime=True&loc=UTC",
    "app.env.data.OBSERVABILITY_LOGGING_LEVEL=info",
    ])
    
k8s_yaml(yaml)
k8s_resource('live-epg', port_forwards=8080)
k8s_resource('live-epg-mariadb', port_forwards=3306)