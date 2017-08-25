# pending_task_rescaler

WIP: very basic feature set for now

If tasks are pending because of missing resources (--reserved-cpu ...), swarm does only start these tasks if new nodes with enough resources get added.
This service also tries to start pending tasks if running containers get stopped (and free some allocated resources) on already used swarm nodes.

start with:

`docker service create --mode global --name rescaler --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock -e MANAGER_URL=tcp://<SWARM_MANAGER>:<PORT> planktom/pending_task_rescaler`
