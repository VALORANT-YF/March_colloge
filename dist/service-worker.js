if(!self.define){let s,e={};const l=(l,n)=>(l=new URL(l+".js",n).href,e[l]||new Promise((e=>{if("document"in self){const s=document.createElement("script");s.src=l,s.onload=e,document.head.appendChild(s)}else s=l,importScripts(l),e()})).then((()=>{let s=e[l];if(!s)throw new Error(`Module ${l} didn’t register its module`);return s})));self.define=(n,i)=>{const r=s||("document"in self?document.currentScript.src:"")||location.href;if(e[r])return;let u={};const o=s=>l(s,r),c={module:{uri:r},exports:u,require:o};e[r]=Promise.all(n.map((s=>c[s]||o(s)))).then((s=>(i(...s),u)))}}define(["./workbox-5b385ed2"],(function(s){"use strict";s.setCacheNameDetails({prefix:"vue_jsbk"}),self.addEventListener("message",(s=>{s.data&&"SKIP_WAITING"===s.data.type&&self.skipWaiting()})),s.precacheAndRoute([{url:"/css/16.585d67ad.css",revision:null},{url:"/css/254.72366bca.css",revision:null},{url:"/css/283.5b677d47.css",revision:null},{url:"/css/326.5901582f.css",revision:null},{url:"/css/588.47ae5507.css",revision:null},{url:"/css/795.a07b347b.css",revision:null},{url:"/css/app.9afc3594.css",revision:null},{url:"/css/chunk-vendors.10dd4e95.css",revision:null},{url:"/fonts/element-icons.f1a45d74.ttf",revision:null},{url:"/fonts/element-icons.ff18efd1.woff",revision:null},{url:"/index.html",revision:"5d0fb8d364f2648fbee4f8a51e330623"},{url:"/js/16.13681fd8.js",revision:null},{url:"/js/254.1f6fa685.js",revision:null},{url:"/js/283.8cb6c514.js",revision:null},{url:"/js/326.66a18a0c.js",revision:null},{url:"/js/588.48ff715b.js",revision:null},{url:"/js/795.8d5774aa.js",revision:null},{url:"/js/app.3e4a62a8.js",revision:null},{url:"/js/chunk-vendors.6d3f184c.js",revision:null},{url:"/manifest.json",revision:"257bfd47f529ef067a513b893c8be0c1"},{url:"/robots.txt",revision:"b6216d61c03e6ce0c9aea6ca7808f7ca"}],{})}));
//# sourceMappingURL=service-worker.js.map