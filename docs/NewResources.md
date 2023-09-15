# Creating New Resources

Mission Control CRDs depend completely on their resources in Crossplane. Because this can be subject to change and needs to accomodate more than one cloud provider, we need to make our resources as generic as possible while also trying to keep cloud specific properties available. This won't always be possible, but to ensure that this stays true most of the time we will contain the following rules when creating a new cloud resource definition in the form of Kubernetes CRD's.

- All of Mission Controls CRD need a controller, this is done with kubebuilders scaffolding.
- The Mission Control CRD must have cluster scope: Some of Crossplane's CRDs have cluster scope and therefore cannot be owned by an object in namespace scope. To avoid this, simply standardize that all new resources must have cluster scope.
- Ownership of resources created: All objects must have ownership of their respective created resources. This allows kubernetes to update and delete these objects whenever a change occurs in the Mission Control CRD or in Crossplanes CRD.
- Group based on what resource is being provided: To ensure that everything stays organized, ensure to group based on Mission Controls groups. For example, compute for VMs or storage groups for buckets.
- Plan similar resources from different providers ahead of time: Not all of the resources have a clear 1-to-1 pairing between them, so planning ahead will not only make it easier to name the Mission Control CRD and it's intention, but also allow for us to look into common parameters beforhand.

## New API Resource (CRD)

To create a new Mission Control CRD use the following command with kubebuilder

`kubebuilder create api --group <GROUP> --version <GROUP_VERSION> --kind <RESOURCE_NAME> --namespaced false`

Choose yes to Resource and Controller. This will generate the basic code for adding our CRD, but we still need to create and apply our definitions to kubernetes. All of this can be handled by running `make install`, but this can be done later.

### Editing CRD

The default values for the resource status and spec need to be changed. This can be done in the file `api/<GROUP>/<GROUP_VERSION>/<RESOURCE>_types.go`. Also ensure that the following line was correctly added and that the CRD is set to cluster-scoped.

`//+kubebuilder:resource:scope=Cluster`

Once you are finished run the `make install` command to finish your setup. No logic will be implemented just yet but you will be able to create a sample and apply it. Basic sample structure should be stored in `config/samples/<GROUP>_<GROUP_VERSION>_<RESOURCE>.yaml`, where a template should already exist to be edited.

### Visualizing Status or Spec parameters.

Although you can add different Spec and Status definitions, these will not appear when querying an applied resource with `kubectl get <RESOURCE>`

To change this go back to the file we modified before in the api folder, and add the following line:
```
//+kubebuilder:printcolumn:name="COLUMN_NAME",type=string,JSONPath=".spec.<RESOURCE_DATA>"
```

This will modify the table out put on `kubectl get <RESOURCE>` and `.spec.<RESOURCE_DATA>` can be changed to whichever path the data you have is located in. This will also work with changing spec to status and supports different data types than just string.

To output the age of the resource specifically use the following line

```
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
```

## Defining CRD Controller

The previous steps added a "kubernetes object" into the system, that we can now apply and get information from, using kubectl commands. That said, there is no functionality, just stored information.

The logic for a kubernetes resource is inside something called the controller. This should be in the file `internal/controller/<GROUP>/<RESOURCE>_controller.go`.

To document this process would be too long so instead look at the documentation about [reconciling best practices](./ReconcilingStrategies) and refer to the current code.

### Testing the controller

If you are using VSCode, get the testing file from `hack/hoftherose/public/vscode/launch.json` Copy this into `.vscode/launch.json` and you should be able to run a testing evironment now. Note that unless the code is actively running, the CRDs will only serve as information. CRDs do not have any functionality unless paired with a running controller.

### Deploying the controller

In order to run the controller full time you need to create an image and run this image on the kubernetes cluster.
>Note: This hasn't been done yet so documentation is in standby.
