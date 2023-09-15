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

### Creating finalizers

Finalizers are ways to ensure that the changes we make are undone before deleting our Custom Resource. This can be achieved by giving ownership to the object, but finalizers ensure that the order is maintained. Do not delete this CR until I have cleaned up successfully. In implementation, finalizers are just strings added and removed in the CR yaml.

If the list is empty, then there are no further steps to be taken and the CR can be deleted with no problem. Otherwise, we need to review the strings to get an idea of what steps need to be taken. For example: "delete_database_entry" will tell us what we need to know.

Most of the time these strings need to be added and maintained when the CR is first created, best practices is to make the change, then add the string. The same goes for deleting, first cleanup, then remove the finalizer.

CR will NOT be deleted even if it has an expired deleted timestamp until all of its finalizers are removed.

### Adding resource to Scheme
