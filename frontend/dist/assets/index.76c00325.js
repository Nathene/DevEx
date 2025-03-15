(function(){const e=document.createElement("link").relList;if(e&&e.supports&&e.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))i(o);new MutationObserver(o=>{for(const l of o)if(l.type==="childList")for(const s of l.addedNodes)s.tagName==="LINK"&&s.rel==="modulepreload"&&i(s)}).observe(document,{childList:!0,subtree:!0});function n(o){const l={};return o.integrity&&(l.integrity=o.integrity),o.referrerpolicy&&(l.referrerPolicy=o.referrerpolicy),o.crossorigin==="use-credentials"?l.credentials="include":o.crossorigin==="anonymous"?l.credentials="omit":l.credentials="same-origin",l}function i(o){if(o.ep)return;o.ep=!0;const l=n(o);fetch(o.href,l)}})();function C(){}function Z(t){return t()}function J(){return Object.create(null)}function O(t){t.forEach(Z)}function x(t){return typeof t=="function"}function te(t,e){return t!=t?e==e:t!==e||t&&typeof t=="object"||typeof t=="function"}function ne(t){return Object.keys(t).length===0}function v(t,e){t.appendChild(e)}function _(t,e,n){t.insertBefore(e,n||null)}function g(t){t.parentNode&&t.parentNode.removeChild(t)}function se(t,e){for(let n=0;n<t.length;n+=1)t[n]&&t[n].d(e)}function a(t){return document.createElement(t)}function $(t){return document.createTextNode(t)}function k(){return $(" ")}function ie(t,e,n,i){return t.addEventListener(e,n,i),()=>t.removeEventListener(e,n,i)}function c(t,e,n){n==null?t.removeAttribute(e):t.getAttribute(e)!==n&&t.setAttribute(e,n)}function le(t){return Array.from(t.childNodes)}function D(t,e){e=""+e,t.data!==e&&(t.data=e)}let G;function I(t){G=t}function oe(){if(!G)throw new Error("Function called outside component initialization");return G}function re(t){oe().$$.on_mount.push(t)}const E=[],Q=[];let H=[];const W=[],fe=Promise.resolve();let q=!1;function ce(){q||(q=!0,fe.then(ee))}function z(t){H.push(t)}const F=new Set;let T=0;function ee(){if(T!==0)return;const t=G;do{try{for(;T<E.length;){const e=E[T];T++,I(e),de(e.$$)}}catch(e){throw E.length=0,T=0,e}for(I(null),E.length=0,T=0;Q.length;)Q.pop()();for(let e=0;e<H.length;e+=1){const n=H[e];F.has(n)||(F.add(n),n())}H.length=0}while(E.length);for(;W.length;)W.pop()();q=!1,F.clear(),I(t)}function de(t){if(t.fragment!==null){t.update(),O(t.before_update);const e=t.dirty;t.dirty=[-1],t.fragment&&t.fragment.p(t.ctx,e),t.after_update.forEach(z)}}function ae(t){const e=[],n=[];H.forEach(i=>t.indexOf(i)===-1?e.push(i):n.push(i)),n.forEach(i=>i()),H=e}const ue=new Set;function me(t,e){t&&t.i&&(ue.delete(t),t.i(e))}function ve(t,e,n,i){const{fragment:o,after_update:l}=t.$$;o&&o.m(e,n),i||z(()=>{const s=t.$$.on_mount.map(Z).filter(x);t.$$.on_destroy?t.$$.on_destroy.push(...s):O(s),t.$$.on_mount=[]}),l.forEach(z)}function ge(t,e){const n=t.$$;n.fragment!==null&&(ae(n.after_update),O(n.on_destroy),n.fragment&&n.fragment.d(e),n.on_destroy=n.fragment=null,n.ctx=[])}function _e(t,e){t.$$.dirty[0]===-1&&(E.push(t),ce(),t.$$.dirty.fill(0)),t.$$.dirty[e/31|0]|=1<<e%31}function pe(t,e,n,i,o,l,s,r=[-1]){const y=G;I(t);const d=t.$$={fragment:null,ctx:[],props:l,update:C,not_equal:o,bound:J(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(e.context||(y?y.$$.context:[])),callbacks:J(),dirty:r,skip_bound:!1,root:e.target||y.$$.root};s&&s(d.root);let u=!1;if(d.ctx=n?n(t,e.props||{},(p,w,...m)=>{const f=m.length?m[0]:w;return d.ctx&&o(d.ctx[p],d.ctx[p]=f)&&(!d.skip_bound&&d.bound[p]&&d.bound[p](f),u&&_e(t,p)),w}):[],d.update(),u=!0,O(d.before_update),d.fragment=i?i(d.ctx):!1,e.target){if(e.hydrate){const p=le(e.target);d.fragment&&d.fragment.l(p),p.forEach(g)}else d.fragment&&d.fragment.c();e.intro&&me(t.$$.fragment),ve(t,e.target,e.anchor,e.customElement),ee()}I(y)}class he{$destroy(){ge(this,1),this.$destroy=C}$on(e,n){if(!x(n))return C;const i=this.$$.callbacks[e]||(this.$$.callbacks[e]=[]);return i.push(n),()=>{const o=i.indexOf(n);o!==-1&&i.splice(o,1)}}$set(e){this.$$set&&!ne(e)&&(this.$$.skip_bound=!0,this.$$set(e),this.$$.skip_bound=!1)}}function ke(){return window.go.main.App.GetCPUDetails()}function ye(){return window.go.main.App.GetCPUInfo()}function be(){return window.go.main.App.GetDiskDetails()}function we(){return window.go.main.App.GetDiskInfo()}function Le(){return window.go.main.App.GetRAMDetails()}function Ce(){return window.go.main.App.GetRAMInfo()}function X(t,e,n){const i=t.slice();return i[10]=e[n],i}function Y(t){let e,n,i=t[10].icon+"",o,l,s=t[10].name+"",r,y,d,u,p;function w(){return t[8](t[10])}return{c(){e=a("button"),n=a("span"),o=$(i),l=k(),r=$(s),y=k(),c(n,"class","icon svelte-gmnf65"),c(e,"class",d="nav-item "+(t[6]===t[10].id?"active":"")+" svelte-gmnf65")},m(m,f){_(m,e,f),v(e,n),v(n,o),v(e,l),v(e,r),v(e,y),u||(p=ie(e,"click",w),u=!0)},p(m,f){t=m,f&64&&d!==(d="nav-item "+(t[6]===t[10].id?"active":"")+" svelte-gmnf65")&&c(e,"class",d)},d(m){m&&g(e),u=!1,p()}}}function Me(t){let e,n,i,o,l;return{c(){e=a("div"),e.innerHTML=`<h3 class="svelte-gmnf65">App Settings</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,n=k(),i=a("div"),i.innerHTML=`<h3 class="svelte-gmnf65">Theme Settings</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,o=k(),l=a("div"),l.innerHTML=`<h3 class="svelte-gmnf65">Notification Settings</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,c(e,"class","card svelte-gmnf65"),c(i,"class","card svelte-gmnf65"),c(l,"class","card svelte-gmnf65")},m(s,r){_(s,e,r),_(s,n,r),_(s,i,r),_(s,o,r),_(s,l,r)},p:C,d(s){s&&g(e),s&&g(n),s&&g(i),s&&g(o),s&&g(l)}}}function Se(t){let e,n,i,o,l;return{c(){e=a("div"),e.innerHTML=`<h3 class="svelte-gmnf65">System Logs</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,n=k(),i=a("div"),i.innerHTML=`<h3 class="svelte-gmnf65">Application Logs</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,o=k(),l=a("div"),l.innerHTML=`<h3 class="svelte-gmnf65">Error Logs</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,c(e,"class","card svelte-gmnf65"),c(i,"class","card svelte-gmnf65"),c(l,"class","card svelte-gmnf65")},m(s,r){_(s,e,r),_(s,n,r),_(s,i,r),_(s,o,r),_(s,l,r)},p:C,d(s){s&&g(e),s&&g(n),s&&g(i),s&&g(o),s&&g(l)}}}function $e(t){let e,n,i,o,l;return{c(){e=a("div"),e.innerHTML=`<h3 class="svelte-gmnf65">Running Processes</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,n=k(),i=a("div"),i.innerHTML=`<h3 class="svelte-gmnf65">Process Tree</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,o=k(),l=a("div"),l.innerHTML=`<h3 class="svelte-gmnf65">Resource Usage</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,c(e,"class","card svelte-gmnf65"),c(i,"class","card svelte-gmnf65"),c(l,"class","card svelte-gmnf65")},m(s,r){_(s,e,r),_(s,n,r),_(s,i,r),_(s,o,r),_(s,l,r)},p:C,d(s){s&&g(e),s&&g(n),s&&g(i),s&&g(o),s&&g(l)}}}function Ae(t){let e,n,i,o,l;return{c(){e=a("div"),e.innerHTML=`<h3 class="svelte-gmnf65">Network Interfaces</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,n=k(),i=a("div"),i.innerHTML=`<h3 class="svelte-gmnf65">Active Connections</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,o=k(),l=a("div"),l.innerHTML=`<h3 class="svelte-gmnf65">Network Usage</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,c(e,"class","card svelte-gmnf65"),c(i,"class","card svelte-gmnf65"),c(l,"class","card svelte-gmnf65")},m(s,r){_(s,e,r),_(s,n,r),_(s,i,r),_(s,o,r),_(s,l,r)},p:C,d(s){s&&g(e),s&&g(n),s&&g(i),s&&g(o),s&&g(l)}}}function De(t){let e,n,i,o,l;return{c(){e=a("div"),e.innerHTML=`<h3 class="svelte-gmnf65">Docker Containers</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,n=k(),i=a("div"),i.innerHTML=`<h3 class="svelte-gmnf65">Docker Images</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,o=k(),l=a("div"),l.innerHTML=`<h3 class="svelte-gmnf65">Docker Volumes</h3> 
          <div class="metric svelte-gmnf65">Coming Soon</div>`,c(e,"class","card svelte-gmnf65"),c(i,"class","card svelte-gmnf65"),c(l,"class","card svelte-gmnf65")},m(s,r){_(s,e,r),_(s,n,r),_(s,i,r),_(s,o,r),_(s,l,r)},p:C,d(s){s&&g(e),s&&g(n),s&&g(i),s&&g(o),s&&g(l)}}}function Te(t){let e,n,i,o,l,s,r,y,d,u,p,w,m,f,M,h,A,R,S,N,K,P,j,V,U,B;return{c(){e=a("div"),n=a("h3"),n.textContent="CPU Usage",i=k(),o=a("div"),l=$(t[0]),s=k(),r=a("div"),y=$(t[1]),d=k(),u=a("div"),p=a("h3"),p.textContent="Memory Usage",w=k(),m=a("div"),f=$(t[2]),M=k(),h=a("div"),A=$(t[3]),R=k(),S=a("div"),N=a("h3"),N.textContent="Disk Usage",K=k(),P=a("div"),j=$(t[4]),V=k(),U=a("div"),B=$(t[5]),c(n,"class","svelte-gmnf65"),c(o,"class","metric svelte-gmnf65"),c(r,"class","details svelte-gmnf65"),c(e,"class","card svelte-gmnf65"),c(p,"class","svelte-gmnf65"),c(m,"class","metric svelte-gmnf65"),c(h,"class","details svelte-gmnf65"),c(u,"class","card svelte-gmnf65"),c(N,"class","svelte-gmnf65"),c(P,"class","metric svelte-gmnf65"),c(U,"class","details svelte-gmnf65"),c(S,"class","card svelte-gmnf65")},m(b,L){_(b,e,L),v(e,n),v(e,i),v(e,o),v(o,l),v(e,s),v(e,r),v(r,y),_(b,d,L),_(b,u,L),v(u,p),v(u,w),v(u,m),v(m,f),v(u,M),v(u,h),v(h,A),_(b,R,L),_(b,S,L),v(S,N),v(S,K),v(S,P),v(P,j),v(S,V),v(S,U),v(U,B)},p(b,L){L&1&&D(l,b[0]),L&2&&D(y,b[1]),L&4&&D(f,b[2]),L&8&&D(A,b[3]),L&16&&D(j,b[4]),L&32&&D(B,b[5])},d(b){b&&g(e),b&&g(d),b&&g(u),b&&g(R),b&&g(S)}}}function Ee(t){let e,n,i,o,l,s,r,y,d=t[7],u=[];for(let f=0;f<d.length;f+=1)u[f]=Y(X(t,d,f));function p(f,M){if(f[6]==="system")return Te;if(f[6]==="docker")return De;if(f[6]==="network")return Ae;if(f[6]==="processes")return $e;if(f[6]==="logs")return Se;if(f[6]==="settings")return Me}let w=p(t),m=w&&w(t);return{c(){e=a("div"),n=a("nav"),i=a("div"),i.innerHTML='<h2 class="svelte-gmnf65">DevEx</h2>',o=k(),l=a("div");for(let f=0;f<u.length;f+=1)u[f].c();s=k(),r=a("main"),y=a("div"),m&&m.c(),c(i,"class","nav-header svelte-gmnf65"),c(l,"class","nav-items svelte-gmnf65"),c(n,"class","side-nav svelte-gmnf65"),c(y,"class","dashboard-grid svelte-gmnf65"),c(r,"class","main-content svelte-gmnf65"),c(e,"class","app-container svelte-gmnf65")},m(f,M){_(f,e,M),v(e,n),v(n,i),v(n,o),v(n,l);for(let h=0;h<u.length;h+=1)u[h]&&u[h].m(l,null);v(e,s),v(e,r),v(r,y),m&&m.m(y,null)},p(f,[M]){if(M&192){d=f[7];let h;for(h=0;h<d.length;h+=1){const A=X(f,d,h);u[h]?u[h].p(A,M):(u[h]=Y(A),u[h].c(),u[h].m(l,null))}for(;h<u.length;h+=1)u[h].d(1);u.length=d.length}w===(w=p(f))&&m?m.p(f,M):(m&&m.d(1),m=w&&w(f),m&&(m.c(),m.m(y,null)))},i:C,o:C,d(f){f&&g(e),se(u,f),m&&m.d()}}}function He(t,e,n){let i="Loading...",o="Loading...",l="Loading...",s="Loading...",r="Loading...",y="Loading...",d="system";const u=[{id:"system",name:"System",icon:"\u{1F4BB}"},{id:"docker",name:"Docker",icon:"\u{1F433}"},{id:"network",name:"Network",icon:"\u{1F310}"},{id:"processes",name:"Processes",icon:"\u2699\uFE0F"},{id:"logs",name:"Logs",icon:"\u{1F4DD}"},{id:"settings",name:"Settings",icon:"\u2699\uFE0F"}];async function p(){try{n(0,i=await ye()),n(1,o=await ke()),n(2,l=await Ce()),n(3,s=await Le()),n(4,r=await we()),n(5,y=await be())}catch(m){console.error("Error updating metrics:",m)}}return re(()=>{p();const m=setInterval(p,2e3);return()=>clearInterval(m)}),[i,o,l,s,r,y,d,u,m=>n(6,d=m.id)]}class Ie extends he{constructor(e){super(),pe(this,e,He,Ee,te,{})}}new Ie({target:document.getElementById("app")});
