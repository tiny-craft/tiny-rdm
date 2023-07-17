export namespace types {
	
	export class Connection {
	    name: string;
	    group?: string;
	    addr?: string;
	    port?: number;
	    username?: string;
	    password?: string;
	    defaultFilter?: string;
	    keySeparator?: string;
	    connTimeout?: number;
	    execTimeout?: number;
	    markColor?: string;
	    type?: string;
	    connections?: Connection[];
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.group = source["group"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.defaultFilter = source["defaultFilter"];
	        this.keySeparator = source["keySeparator"];
	        this.connTimeout = source["connTimeout"];
	        this.execTimeout = source["execTimeout"];
	        this.markColor = source["markColor"];
	        this.type = source["type"];
	        this.connections = this.convertValues(source["connections"], Connection);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class ConnectionConfig {
	    name: string;
	    group?: string;
	    addr?: string;
	    port?: number;
	    username?: string;
	    password?: string;
	    defaultFilter?: string;
	    keySeparator?: string;
	    connTimeout?: number;
	    execTimeout?: number;
	    markColor?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.group = source["group"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.defaultFilter = source["defaultFilter"];
	        this.keySeparator = source["keySeparator"];
	        this.connTimeout = source["connTimeout"];
	        this.execTimeout = source["execTimeout"];
	        this.markColor = source["markColor"];
	    }
	}
	export class JSResp {
	    success: boolean;
	    msg: string;
	    data?: any;
	
	    static createFrom(source: any = {}) {
	        return new JSResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}

}

