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

export interface Warning {
    timestamp: string
    message: string
}



export class ConciergeApi {
    private basePath: string
    constructor(basePath: string) {
        this.basePath = basePath
    }

    public async createSingleDeployment(deployment: SingleContainerDto) {
        await this.$axios.$post()
    }
