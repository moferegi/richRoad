import { baseUrl }from "@/utils/request.js"
export const getUrl = (url) =>{
	if(!url) return "http://localhost:8888"
	const path = baseUrl
	if(url.slice(0, 4) === 'data'){
		return url
	}
	 if (url && url.slice(0, 4) !== 'http'){
	    if (path === "/"){
	        return url
	    }
	    if (url.slice(0, 1) === "/"){
	       return path + url
	    }
	    return path + "/" + url
	  }else{
	    return url
	  }
}