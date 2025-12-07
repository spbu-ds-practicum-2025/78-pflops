include .env

.PHONY: remove-volumes up down

up:
	docker-compose up -d --build

down: 
	docker-compose down

remove-volumes:
	docker volume rm services_ad_service_db_data services_media_service_minio_data services_user_service_db_data
