.PHONY: rabbit
rabbit:
	docker run -d --name rabbit -p 5672:5672 -p 15672:15672 rabbitmq:management


rabbit_existed:
	docker start  rabbit


rabbit_add_user: rabbit_create_user rabbit_add_user_perm

rabbit_create_user:
	docker exec rabbit rabbitmqctl add_user gorik 1
rabbit_add_user_perm:
	docker exec rabbit rabbitmqctl set_user_tags gorik administrator\

rabbit_add_vhost:
	docker exec rabbit rabbitmqctl add_vhost army



rabbit_set_permissions:
	docker exec rabbit rabbitmqctl set_permissions -p army gorik ".*" ".*" ".*"


rabbit_echange:
	docker exec rabbit rabbitmqadmin declare  exchange name=army_event type=fanout durable=true --vhost=army -u  gorik -p 1

rabbit_del_echange:
	docker exec rabbit rabbitmqadmin delete  exchange name=army_event --vhost=army -u  gorik -p 1


rabbit_perm_echange:
	docker exec rabbit rabbitmqctl set_topic_permissions -p army gorik army_event ".*" ".*"

