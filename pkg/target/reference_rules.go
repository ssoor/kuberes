package target

const ruls = `
- kind: Pod
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/volumes/secret/secretName
        - spec/containers/env/valueFrom/secretKeyRef/name
        - spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/containers/envFrom/secretRef/name
        - spec/initContainers/envFrom/secretRef/name
        - spec/imagePullSecrets/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/volumes/configMap/name
        - spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/containers/envFrom/configMapRef/name
        - spec/initContainers/envFrom/configMapRef/name
    - kind: ServiceAccount
      version: v1
      paths: 
        - spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/volumes/persistentVolumeClaim/claimName

- kind: Deployment
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/template/spec/volumes/secret/secretName
        - spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/containers/envFrom/secretRef/name
        - spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/template/spec/imagePullSecrets/name
        - spec/template/spec/volumes/projected/sources/secret/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/template/spec/volumes/configMap/name
        - spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/containers/envFrom/configMapRef/name
        - spec/template/spec/initContainers/envFrom/configMapRef/name
        - spec/template/spec/volumes/projected/sources/configMap/name
    - kind: ServiceAccount
      version: v1
      paths: 
        - spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/selector/matchLabels
        - spec/template/metadata/labels
        - spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
        - spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
        - spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
        - spec/template/spec/affinity/podAntiAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: StatefulSet
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/template/spec/volumes/secret/secretName
        - spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/containers/envFrom/secretRef/name
        - spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/template/spec/imagePullSecrets/name
        - spec/template/spec/volumes/projected/sources/secret/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/template/spec/volumes/configMap/name
        - spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/containers/envFrom/configMapRef/name
        - spec/template/spec/initContainers/envFrom/configMapRef/name
        - spec/template/spec/volumes/projected/sources/configMap/name
    - kind: Service
      group: apps
      paths: 
        - spec/serviceName
    - kind: ServiceAccount
      paths: 
        - spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/selector/matchLabels
        - spec/template/metadata/labels
        - spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
        - spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
        - spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
        - spec/template/spec/affinity/podAntiAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
        - spec/volumeClaimTemplates/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: DaemonSet
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/template/spec/volumes/secret/secretName
        - spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/containers/envFrom/secretRef/name
        - spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/template/spec/imagePullSecrets/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/template/spec/volumes/configMap/name
        - spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/containers/envFrom/configMapRef/name
        - spec/template/spec/initContainers/envFrom/configMapRef/name
    - kind: ServiceAccount
      version: v1
      paths: 
        - spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/selector/matchLabels
        - spec/template/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: ReplicaSet
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/template/spec/volumes/secret/secretName
        - spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/containers/envFrom/secretRef/name
        - spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/template/spec/imagePullSecrets/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/template/spec/volumes/configMap/name
        - spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/containers/envFrom/configMapRef/name
        - spec/template/spec/initContainers/envFrom/configMapRef/name
  metadata.labels:
    - create: true
      paths: 
        - spec/selector/matchLabels
        - spec/template/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: Job
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/template/spec/volumes/secret/secretName
        - spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/template/spec/containers/envFrom/secretRef/name
        - spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/template/spec/imagePullSecrets/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/template/spec/volumes/configMap/name
        - spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/template/spec/containers/envFrom/configMapRef/name
        - spec/template/spec/initContainers/envFrom/configMapRef/name
    - kind: ServiceAccount
      paths: 
        - spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/selector/matchLabels
        - spec/template/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: CronJob
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/jobTemplate/spec/template/spec/volumes/secret/secretName
        - spec/jobTemplate/spec/template/spec/containers/env/valueFrom/secretKeyRef/name
        - spec/jobTemplate/spec/template/spec/initContainers/env/valueFrom/secretKeyRef/name
        - spec/jobTemplate/spec/template/spec/containers/envFrom/secretRef/name
        - spec/jobTemplate/spec/template/spec/initContainers/envFrom/secretRef/name
        - spec/jobTemplate/spec/template/spec/imagePullSecrets/name
    - kind: ConfigMap
      version: v1
      paths: 
        - spec/jobTemplate/spec/template/spec/volumes/configMap/name
        - spec/jobTemplate/spec/template/spec/containers/env/valueFrom/configMapKeyRef/name
        - spec/jobTemplate/spec/template/spec/initContainers/env/valueFrom/configMapKeyRef/name
        - spec/jobTemplate/spec/template/spec/containers/envFrom/configMapRef/name
        - spec/jobTemplate/spec/template/spec/initContainers/envFrom/configmapRef/name
    - kind: ServiceAccount
      paths: 
        - spec/jobTemplate/spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      version: v1
      paths: 
        - spec/jobTemplate/spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/jobTemplate/spec/selector/matchLabels
        - spec/jobTemplate/metadata/labels
        - spec/jobTemplate/spec/template/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/jobTemplate/metadata/annotations
        - spec/jobTemplate/spec/template/metadata/annotations

- kind: Ingress
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - spec/tls/secretName
        - metadata/annotations/ingress.kubernetes.io\/auth-secret
        - metadata/annotations/nginx.ingress.kubernetes.io\/auth-secret
    - kind: Service
      version: v1
      paths: 
        - spec/backend/serviceName
        - spec/rules/http/paths/backend/serviceName

- kind: StorageClass
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - parameters/secretName
        - parameters/adminSecretName
        - parameters/userSecretName
        - parameters/secretRef

- kind: Service
  version: v1
  metadata.annotations:
    - create: true
      paths: 
        - spec/selector

- kind: ServiceAccount
  version: v1
  matedata.name:
    - kind: Secret
      version: v1
      paths: 
        - imagePullSecrets/name
          
- kind: HorizontalPodAutoscaler
  version: v1
  matedata.name:
    - kind: ReplicaSet
      paths: 
        - spec/scaleTargetRef/name
    - kind: Deployment
      paths: 
        - spec/scaleTargetRef/name
    - kind: ReplicationController
      paths: 
        - spec/scaleTargetRef/name
        
- kind: ReplicationController
  version: v1
  matedata.name:
    - kind: Role
      group: rbac.authorization.k8s.io
      paths: 
        - roleRef/name
    - kind: ClusterRole
      paths: 
        - roleRef/name
    - kind: ServiceAccount
      paths: 
        - subjects/name
        - spec/template/spec/serviceAccountName
    - kind: PersistentVolumeClaim
      paths: 
        - spec/template/spec/volumes/persistentVolumeClaim/claimName
  metadata.labels:
    - create: true
      paths: 
        - spec/selector
        - spec/template/metadata/labels
  metadata.annotations:
    - create: true
      paths: 
        - spec/template/metadata/annotations

- kind: NetworkPolicy
  version: v1
  metadata.labels:
    - create: true
      paths: 
        - spec/podSelector/matchLabels
        - spec/ingress/from/podSelector/matchLabels
        - spec/egress/to/podSelector/matchLabels

- kind: PodDisruptionBudget
  version: v1
  metadata.labels:
    - create: false
      paths: 
        - spec/selector/matchLabels

- kind: PersistentVolumeClaim
  version: v1
  matedata.name:
    - kind: PersistentVolume
      version: v1
      paths: 
        - spec/volumeName

- kind: Role
  version: v1
  matedata.name:
    - kind: Secret
      paths: 
        - rules/resourceNames
     
- kind: ClusterRole
  version: v1
  matedata.name:
    - kind: Secret
      paths: 
        - rules/resourceNames
     
- kind: RoleBinding
  version: v1
  matedata.name:
    - kind: Role
      group: rbac.authorization.k8s.io
      paths: 
        - roleRef/name
    - kind: ClusterRole
      paths: 
        - roleRef/name
    - kind: ServiceAccount
      paths: 
        - subjects/name

- kind: ClusterRoleBinding
  version: v1
  matedata.name:
    - kind: ClusterRole
      group: rbac.authorization.k8s.io
      paths: 
        - roleRef/name
    - kind: ServiceAccount
      group: rbac.authorization.k8s.io
      paths: 
        - subjects/name
`

// GetReferenceRules is
func GetReferenceRules() []byte {
	return []byte(ruls)
}
