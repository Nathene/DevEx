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

export namespace history {
	
	export class DockerMetrics {
	    // Go type: time
	    timestamp: any;
	    daemonRunning: boolean;
	    containersRunning: number;
	    containersTotal: number;
	    imagesCount: number;
	
	    static createFrom(source: any = {}) {
	        return new DockerMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.daemonRunning = source["daemonRunning"];
	        this.containersRunning = source["containersRunning"];
	        this.containersTotal = source["containersTotal"];
	        this.imagesCount = source["imagesCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NetworkMetrics {
	    // Go type: time
	    timestamp: any;
	    internetConnected: boolean;
	    pingLatency: number;
	    dnsWorking: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.internetConnected = source["internetConnected"];
	        this.pingLatency = source["pingLatency"];
	        this.dnsWorking = source["dnsWorking"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TimeSeriesPoint {
	    // Go type: time
	    timestamp: any;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new TimeSeriesPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.value = source["value"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
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

export namespace process {
	
	export class PortInfo {
	    port: number;
	    protocol: string;
	    localAddr: string;
	    state: string;
	    pid: number;
	
	    static createFrom(source: any = {}) {
	        return new PortInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.protocol = source["protocol"];
	        this.localAddr = source["localAddr"];
	        this.state = source["state"];
	        this.pid = source["pid"];
	    }
	}
	export class ProcessInfo {
	    pid: number;
	    name: string;
	    commandLine: string;
	    username: string;
	    cpuPercent: number;
	    memoryUsage: number;
	    // Go type: time
	    startTime: any;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.commandLine = source["commandLine"];
	        this.username = source["username"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memoryUsage = source["memoryUsage"];
	        this.startTime = this.convertValues(source["startTime"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProcessWithPorts {
	    process: ProcessInfo;
	    ports: PortInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ProcessWithPorts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.process = this.convertValues(source["process"], ProcessInfo);
	        this.ports = this.convertValues(source["ports"], PortInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

