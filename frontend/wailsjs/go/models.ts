export namespace config {
	
	export class WindowConfig {
	    width: number;
	    height: number;
	    minWidth: number;
	    minHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.minWidth = source["minWidth"];
	        this.minHeight = source["minHeight"];
	    }
	}
	export class AppConfig {
	    productName: string;
	    productNameEn: string;
	    version: string;
	    description: string;
	    author: Record<string, string>;
	    website: string;
	    sponsorQrCode: string;
	    defaultWindow: WindowConfig;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productName = source["productName"];
	        this.productNameEn = source["productNameEn"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.website = source["website"];
	        this.sponsorQrCode = source["sponsorQrCode"];
	        this.defaultWindow = this.convertValues(source["defaultWindow"], WindowConfig);
	        this.language = source["language"];
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

export namespace hosts {
	
	export class BackupInfo {
	    id: string;
	    // Go type: time
	    timestamp: any;
	    size: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.size = source["size"];
	        this.path = source["path"];
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
	export class HostEntry {
	    ip: string;
	    domain: string;
	    comment: string;
	    alternate_ip: string;
	    description: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HostEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.domain = source["domain"];
	        this.comment = source["comment"];
	        this.alternate_ip = source["alternate_ip"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	    }
	}
	export class HostsProfile {
	    name: string;
	    entries: HostEntry[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new HostsProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.entries = this.convertValues(source["entries"], HostEntry);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

export namespace nvm {
	
	export class NodeVersion {
	    version: string;
	    path: string;
	    isCurrent: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NodeVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.path = source["path"];
	        this.isCurrent = source["isCurrent"];
	    }
	}

}

