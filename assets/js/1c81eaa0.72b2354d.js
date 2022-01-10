"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[6034],{3905:function(e,n,t){t.d(n,{Zo:function(){return l},kt:function(){return g}});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function c(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var d=r.createContext({}),u=function(e){var n=r.useContext(d),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},l=function(e){var n=u(e.components);return r.createElement(d.Provider,{value:n},e.children)},s={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},p=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,i=e.originalType,d=e.parentName,l=c(e,["components","mdxType","originalType","parentName"]),p=u(t),g=o,f=p["".concat(d,".").concat(g)]||p[g]||s[g]||i;return t?r.createElement(f,a(a({ref:n},l),{},{components:t})):r.createElement(f,a({ref:n},l))}));function g(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=p;var c={};for(var d in n)hasOwnProperty.call(n,d)&&(c[d]=n[d]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var u=2;u<i;u++)a[u]=t[u];return r.createElement.apply(null,a)}return r.createElement.apply(null,t)}p.displayName="MDXCreateElement"},613:function(e,n,t){t.r(n),t.d(n,{frontMatter:function(){return c},contentTitle:function(){return d},metadata:function(){return u},toc:function(){return l},default:function(){return p}});var r=t(7462),o=t(3366),i=(t(7294),t(3905)),a=["components"],c={id:"node-config",title:"Query a Node's Configuration"},d=void 0,u={unversionedId:"getting-started/queries/curl/node-config",id:"getting-started/queries/curl/node-config",isDocsHomePage:!1,title:"Query a Node's Configuration",description:"Querying a Particular Node's Configuration",source:"@site/docs/getting-started/queries/curl/node-config.md",sourceDirName:"getting-started/queries/curl",slug:"/getting-started/queries/curl/node-config",permalink:"/orion-server/docs/getting-started/queries/curl/node-config",tags:[],version:"current",frontMatter:{id:"node-config",title:"Query a Node's Configuration"}},l=[{value:"Querying a Particular Node&#39;s Configuration",id:"querying-a-particular-nodes-configuration",children:[],level:3}],s={toc:l};function p(e){var n=e.components,t=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},s,t,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h3",{id:"querying-a-particular-nodes-configuration"},"Querying a Particular Node's Configuration"),(0,i.kt)("p",null,"While ",(0,i.kt)("inlineCode",{parentName:"p"},"GET /config/tx")," returns the complete configuration, we can use ",(0,i.kt)("inlineCode",{parentName:"p"},"GET /config/node/{nodeid}")," to fetch the configuration of a particular node.\nHere, the submitting user needs to sign ",(0,i.kt)("inlineCode",{parentName:"p"},'{"user_id":"<userid>","node_id":"<nodeid>"}'),", where ",(0,i.kt)("inlineCode",{parentName:"p"},"<userid>")," denotes the id of the submitting user and\n",(0,i.kt)("inlineCode",{parentName:"p"},"<nodeid>")," denotes the id of the blockchain database for which the user is trying to fetch the configuration."),(0,i.kt)("p",null,"In our example, the JSON data to be signed is ",(0,i.kt)("inlineCode",{parentName:"p"},'{"user_id":"admin","node_id":"bdb-node-1"}'),", because the ",(0,i.kt)("inlineCode",{parentName:"p"},"admin")," is submitting a request to fetch\nthe configuration of the node with the id ",(0,i.kt)("inlineCode",{parentName:"p"},"bdb-node-1"),"."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-shell"},'./bin/signer -privatekey=deployment/sample/crypto/admin/admin.key -data=\'{"user_id":"admin","node_id":"bdb-node-1"}\'\n')),(0,i.kt)("p",null,"The above command produces a digital signature and prints it as a base64-encoded string as shown below."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"MEQCIEAxqk1RnWBvCt9OV7zycK3U5itkEumSpEXivv/2R030AiA5mZMgvIDEgPP2J1TGRPsgyXZY6VrAcWTyFAzlLgKpPg==\n")),(0,i.kt)("p",null,"Once the signature is computed, we can issue a ",(0,i.kt)("inlineCode",{parentName:"p"},"GET")," request using the following ",(0,i.kt)("inlineCode",{parentName:"p"},"cURL")," command\nby setting the above signature string in the ",(0,i.kt)("inlineCode",{parentName:"p"},"Signature")," header."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-shell"},'curl \\\n   -H "Content-Type: application/json" \\\n   -H "UserID: admin" \\\n   -H "Signature: MEQCIEAxqk1RnWBvCt9OV7zycK3U5itkEumSpEXivv/2R030AiA5mZMgvIDEgPP2J1TGRPsgyXZY6VrAcWTyFAzlLgKpPg==" \\\n   -X GET http://127.0.0.1:6001/config/node/bdb-node-1 | jq .\n')),(0,i.kt)("p",null,"A sample output of above command is shown below. The actual content might change depending on the configuration specified in ",(0,i.kt)("inlineCode",{parentName:"p"},"config.yml")," for the node ",(0,i.kt)("inlineCode",{parentName:"p"},"bdb-node-1"),". The output would only contain the configuration of the node ",(0,i.kt)("inlineCode",{parentName:"p"},"bdb-node-1"),"."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-webmanifest"},'{\n  "response": {\n    "header": {\n      "node_id": "bdb-node-1"\n    },\n    "node_config": {\n      "id": "bdb-node-1",\n      "address": "127.0.0.1",\n      "port": 6001,\n      "certificate": "MIIBsjCCAVigAwIBAgIQYy4vf2+6qRtczxIkNb9fxjAKBggqhkjOPQQDAjAeMRwwGgYDVQQDExNDYXIgcmVnaXN0cnkgUm9vdENBMB4XDTIxMDYxNjExMTMyN1oXDTIyMDYxNjExMTgyN1owJTEjMCEGA1UEAxMaQ2FyIHJlZ2lzdHJ5IENsaWVudCBzZXJ2ZXIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARWWd8GA/dbkxTiNP7x/LoyAc85sfsIxcmiX+4nzzB3s4SXA+N8YMSiKpsi6gSOVkxKLu43lla2ajyL0Z4WqWvYo3EwbzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAU7nVzp7gto++BPlj5KAF1IA62TNEwDwYDVR0RBAgwBocEfwAAATAKBggqhkjOPQQDAgNIADBFAiEAuyIC0jlV/KcyB8Cz2p3W4aojh+fDCeeRenMwvyP+EcACIDOjDiObMnb/2q2ceAKROr/rzJmakdjkNmw8A0bYL6Pb"\n    }\n  },\n  "signature": "MEUCIB0gWj3xTZ4TkH/FcLhy8X3gVy5n5nOOIz1949HtPqhbAiEAxRgdd8z2LML0zk++rVNldJO+VpVH6E2j6lQrM6lboME="\n}\n')))}p.isMDXComponent=!0}}]);