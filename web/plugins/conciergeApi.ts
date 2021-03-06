export interface SingleContainerDto {
  containerImage: string
  cpu: number
  ram: number
  gpu: number
  volumeSize: number
  volumeMountPath: string
}
export interface Status {
  name: string
}
export interface Volume {
  size: number
}

export interface Deployment {
  id: number
  name: string
  name_k8s: string
  cpu: number
  ram: number
  gpu: number
  volume: Volume
  status: Status
}

export interface DistributedDeployment {
  id: number
  name: string
  name_k8s: string
  container_image: string
  worker_count: number
  launcher_cpu: number
  launcher_ram: number
  worker_cpu: number
  worker_ram: number
  worker_gpu: number
  status: Status
}

export interface Warning {
  timestamp: string
  reason: string
  message: string
}

export interface Bucket {
  id: number
  name: string
  status: Status
}
