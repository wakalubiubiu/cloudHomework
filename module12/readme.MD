kubectl create ns istio-test
kubectl label ns istio-test istio-injection=enabled
kubectl apply -f service0.yaml -n istio-test
kubectl apply -f service1.yaml -n istio-test
kubectl apply -f service2.yaml -n istio-test
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/CN=*.istio.test' -keyout istio.test.key -out
kubectl create -n istio-system secret tls istio-test --key=istio.test.key --cert=istio.test.crt
curl --resolve httpserver.istio.test:443:$INGRESS_IP https://httpserver.istio.test/service0 -v -k
