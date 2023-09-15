## Reconciling Best Practices

The following points give us

### Ownership of created resources

This requires two steps to be taken. First, the CRD must be assigned ownership to whatever resource it creates. This needs to be done explicitly for each resource. Below is an example:

```
if err := controllerutil.SetControllerReference(CRD, &resource, r.Scheme); err != nil {
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

### Adding resource to Scheme
