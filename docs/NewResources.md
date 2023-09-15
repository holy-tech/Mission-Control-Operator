# Creating New Resources

Mission Control CRDs depend completely on their resources in Crossplane. Because this can be subject to change and needs to accomodate more than one cloud provider, we need to make our resources as generic as possible while also trying to keep cloud specific properties available. This won't always be possible, but to ensure that this stays true most of the time we will contain the following rules when creating a new cloud resource definition in the form of Kubernetes CRD's.

- The Mission Control CRD must have cluster scope: Some of Crossplane's CRDs have cluster scope and therefore cannot be owned by an object in namespace scope. To avoid this, simply standardize that all new resources must have cluster scope.
- Ownership of resources created: All objects must have ownership of their respective created resources. This allows kubernetes to update and delete these objects whenever a change occurs in the Mission Control CRD or in Crossplanes CRD.
- Group based on what resource is being provided: To ensure that everything stays organized, ensure to group based on Mission Controls groups. For example, compute for VMs or storage groups for buckets.
- Plan similar resources from different providers ahead of time: Not all of the resources have a clear 1-to-1 pairing between them, so planning ahead will not only make it easier to name the Mission Control CRD and it's intention, but also allow for us to look into common parameters beforhand.

## New API Resource

To create a new Mission Control CRD use the following command with kubebuilder

`kubebuilder create api --group <GROUP> --version <GROUP_VERSION> --kind <RESOURCE_NAME> --namespaced false`

