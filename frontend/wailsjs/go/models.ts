export namespace main {
	
	export class MyThing {
	    time_str: string;
	
	    static createFrom(source: any = {}) {
	        return new MyThing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time_str = source["time_str"];
	    }
	}

}

