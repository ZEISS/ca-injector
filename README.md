# ca-injector

FOR DEVELOPMENT AND TESTING PURPOSES ONLY

This is a kubernetes [MutatingAdmissionWebhook][mutating_admission_webhook_url] to inject
certificate bundles into pods based on labels and annotations.
With that off-the-shelf containers can be deployed in clusters with custom certificate
authorities, with minimal disruption and minimal maintenance. No more creating images from
upstream base images just to `ADD custom-ca.crt /usr/local/share/ca-certificates/` and
`RUN update-ca-certificates` etc.

[mutating_admission_webhook_url]: https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/

This webhook will be triggered in following cases:
- object label `ca-injector.zeiss.com/inject: "true"`
- namespace label `ca-injector.zeiss.com/inject: "true"` and
  object label `ca-injector.zeiss.com/inject` does not exist
- deployment setting `admissionWebhook.enableNamespacesByDefault: true` (helm chart) and
  namespace/object label `ca-injector.zeiss.com/inject` does not exist

If triggered, this webhook does three things:
1. Add a volume to the pod referencing a certificate bundle specified by the value of the object annotation
   `ca-injector.zeiss.com/inject-ca-from` or the deployment setting `caBundle.configMap` (helm chart).
   The value should correspond with a config map in the same namespace as the pod which has the key `ca.crt`
   containing the CA bundle content.
1. Add this volume to all containers and init containers as a volume mount
1. Add the `SSL_CERT_FILE` and `NODE_EXTRA_CA_CERTS` environment variable [respected by
   OpenSSL](https://www.openssl.org/docs/man3.1/man3/SSL_CTX_set_default_verify_paths.html)
   and most tls libraries.

It is strongly recommended to use this webhook with
[replicator](https://github.com/mittwald/kubernetes-replicator) or [trust-manager](https://github.com/cert-manager/trust-manager) for a consistent experience across namespaces.

## Usage

[Helm](https://helm.sh) must be installed to use the chart. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  `helm repo add ca-injector https://zeiss.github.io/ca-injector`

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo ca-injector`
to see the charts.

To install the ca-injector chart:

  `helm install ca-injector ca-injector/ca-injector`

To uninstall the chart:

  `helm delete ca-injector`

## Example

#TODO
Deploy this in your cluster (see Usage chapter) and create a CA bundle as e.g. `foo-crt` config map, with the key `ca.crt`:

```bash
kubectl create configmap foo-crt --from-file=ca.crt=my-bundle.crt
```

Use this CA bundle by defining the label `ca-injector.zeiss.com/inject:true` and
the annotation `ca-injector.zeiss.com/inject-ca-from: foo-crt` on your pod or
in your helm chart's appropriate annotations section.
`ca.crt` can be changed by configuration `caBundle.crt` in any of the typical
ways (config files at `/etc/ca-injector.yaml`, `$HOME/.config/ca-injector.yaml`,
or environment variable `CAINJECTOR_CABUNDLE_CRT`).


## Release

### App
To trigger a new tagged docker build, create a PR with label 'app-release'. The app Version within the helm chart will be used as reference for the container tag. 
This will be done automatically by below mentioned workflow.

### Helm
In case the appVersion is increased, the helm Chart version should also be increased.
In case the helm Chart version is increased, the appVersion does not have to be increased as well.


Option 1: 
Manually set version and/or appVersion within Helm Chart. The Helm release workflow will create a new release in case the helm Chart version has changed.

Option 2:
Add one or two(app and helm) of the following labels to your PR:
- app-major
- app-minor
- app-patch
- helm-major
- helm-minor
- helm-patch

According to the label, appVersion and/or helm version will be bumped and a PullRequest will be created. The Pull request will include label 'app-release' to trigger above mentioned workflow. After this PR has benn closed, the Helm release workflow will create a new release in case the helm Chart version has changed.

