## Reconciling Best Practices

The following points give us

### Ownership of created resources

This requires two steps to be taken. First, the CustomResource must be assigned ownership to whatever resource it creates. This needs to be done explicitly for each resource. Below is an example:

```
if err := controllerutil.SetControllerReference(CustomResource, &resource, r.Scheme); err != nil {
    return ctrl.Result{}, err
}
```

The second step is to add a watch on all resources the CRD has ownership to. The shortcut for this is with the .Owns method. Below we can see tat MissionKey will watch owned resources that are Secrets or ServicesAccounts.

```
func (r *MissionKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.MissionKey{}).
		Owns(&v1.Secret{}).
		Owns(&v1.ServiceAccount{}).
		Complete(r)
}
```
### Updating created resources

In most cases, changes in the Mission Control CRD yaml needs to be propagated to some of the Crossplane resources owned by it. They also expect the state of these yaml files to stay consistent with what is written in Mission Control, or in other words, we don't want people manually changing the Crossplane objects without going through Mission Control.

There are different ways to update depending on what you changed, but for the most part the steps are:
- Get the object into an empty declaration.
- Make changes to reflect desired outcome.
- Use the r.Update method to upload changes back to kubernetes.

### Creating finalizers

Finalizers are ways to ensure that the changes we make are undone before deleting our Custom Resource. This can be achieved by giving ownership to the object, but finalizers ensure that the order is maintained. Do not delete this CR until I have cleaned up successfully. In implementation, finalizers are just strings added and removed in the CR yaml.

If the list is empty, then there are no further steps to be taken and the CR can be deleted with no problem. Otherwise, we need to review the strings to get an idea of what steps need to be taken. For example: "delete_database_entry" will tell us what we need to know.

Most of the time these strings need to be added and maintained when the CR is first created, best practices is to make the change, then add the string. The same goes for deleting, first cleanup, then remove the finalizer.

CR will NOT be deleted even if it has an expired deleted timestamp until all of its finalizers are removed.

### Adding resource to Scheme

Sometimes we need external resources from seperate operators, for example crossplane and crossplanes providers. If used directly, the items will look for their definitions with the wrong group and version. To fix this, you will need to add the object to the Scheme in the `cmd/main.go` file.

First import the code to the main file, the follow the below structure to create a new SchemeBuilder object.

```
import (
    ...
    cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
    ...
)
func init() {
    ...
    crossplaneSchemeBuilder := &controllerscheme.Builder{GroupVersion: apischeme.GroupVersion{Group: "pkg.crossplane.io", Version: "v1"}}
    crossplaneSchemeBuilder.Register(
        &cpv1.Provider{},
        &cpv1.ProviderList{},
    )
    if err := crossplaneSchemeBuilder.AddToScheme(scheme); err != nil {
        os.Exit(1)
    }
    ...
}
```
