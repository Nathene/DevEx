export namespace docker {
	
	export class Metrics {
	    imagesCount: number;
	    containersAll: number;
	    containersUp: number;
	    diskUsage: string;
	    networkStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new Metrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imagesCount = source["imagesCount"];
	        this.containersAll = source["containersAll"];
	        this.containersUp = source["containersUp"];
	        this.diskUsage = source["diskUsage"];
	        this.networkStatus = source["networkStatus"];
	    }
	}
	export class Status {
	    daemonRunning: boolean;
	    version: string;
	    info: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.daemonRunning = source["daemonRunning"];
	        this.version = source["version"];
	        this.info = source["info"];
	    }
	}

}

export namespace network {
	
	export class Status {
	    internetConnected: boolean;
	    pingLatency: number;
	    pingStatus: string;
	    dnsStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.internetConnected = source["internetConnected"];
	        this.pingLatency = source["pingLatency"];
	        this.pingStatus = source["pingStatus"];
	        this.dnsStatus = source["dnsStatus"];
	    }
	}

}

