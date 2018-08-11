# Book-List

### REST
Built in GO (main.go) for managing booklist. It Supports operations -->
1. add a book
2. show the full booklist
3. delete a book
4. update information about a book

`web` is the binary file of this api.

### Docker

1. write the dockerfile
2. docker build -t kfoozminus/booklistgo:latest .
3. docker run -d -p 5000:8080 --name book_container kfoozminus/booklistgo --rm
    now go to localhost:5000
    Or,
        docker run -d -P kfoozminus/booklistgo
        docker port book_container
        80/tcp -> 0.0.0.0:32769
        443/tcp -> 0.0.0.0:32768
        now go to localhost:32769
4. docker login
5. docker push
6. Remove all docker containers
docker rm $(docker ps -a -q -f status=exited)
7. docker run -it imagename
-it flags attaches us to an interactive tty in the container.

### Kubernetes/Minikube

1. minikube start
2. kubectl run book-kube --image=kfoozminus/booklistgo:v1 --port=8080
3. kubectl get deployments
4. kubectl get pods
5. kubectl expose deployment book-kube --type=LoadBalancer
6. kubectl get services
7. minikube service go-kube --url

Scaling:

8. kubectl scale deployments/book-kube --replicas=3
9. kubectl get pods

Updating:

10. kubectl set image deployment/book-kube booklistgo=booklistgo:latest
11. minikube service book-kube
