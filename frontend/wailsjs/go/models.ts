export namespace log_viewer {
	
	export class ChunkResult {
	    lines: string[];
	    nextStart: number;
	    eof: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ChunkResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lines = source["lines"];
	        this.nextStart = source["nextStart"];
	        this.eof = source["eof"];
	    }
	}
	export class LogFileItem {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new LogFileItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}

export namespace project {
	
	export class ProjectConfig {
	    title: string;
	    key: string;
	    type: string;
	    start: string;
	    stop: string;
	    install: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.key = source["key"];
	        this.type = source["type"];
	        this.start = source["start"];
	        this.stop = source["stop"];
	        this.install = source["install"];
	    }
	}
	export class Project {
	    Config: ProjectConfig;
	    path: string;
	    running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Config = this.convertValues(source["Config"], ProjectConfig);
	        this.path = source["path"];
	        this.running = source["running"];
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

