import{P as SERVER_PATH,O as request,i as useRouter,r as reactive,C as resolveComponent,o as openBlock,c as createElementBlock,a as createBaseVNode,h as createVNode,w as withCtx,t as toDisplayString,u as unref,j as withDirectives,v as vModelText,$ as withKeys,F as Fragment,b as renderList,k as createCommentVNode,d as utils,R as RESPONSE,g as getCurrentInstance,s as watch,V as nextTick,p as normalizeClass,z as createBlock,a0 as IM_ERRORS,a1 as METHOD_MAP,a2 as SIGNAL_TYPE}from"./index-a6509667.js";function getConns(m){let{count:p=50,start:d,app_key:c,user_id:o}=m,i=`${SERVER_PATH.CONN_GET_LIST}?app_key=${c}&user_id=${o}&count=${p}&start=${d}`;return request(i,{method:"GET"})}function getConn(m){let{count:p=50,start:d,app_key:c,session:o}=m,i=`${SERVER_PATH.CONN_GET_ONE}?app_key=${c}&session=${o}&count=${p}&start=${d}`;return request(i,{method:"GET"})}const Inspect={getConns,getConn},_hoisted_1$2={class:"jg-tconmanger-box jg-table-box"},_hoisted_2$2={class:"jg-table-header jg-log-header"},_hoisted_3$2={class:"jg-table-header-lf-box"},_hoisted_4$2={class:"jg-table-lf-item"},_hoisted_5$2=["onClick"],_hoisted_6$2={class:"jg-table-lf-item"},_hoisted_7$2={class:"jg-table-body"},_hoisted_8$2={class:"table jg-table"},_hoisted_9$1={class:"jg-td-c"},_hoisted_10$1={class:"jg-td-c"},_hoisted_11$1={class:"jg-td-c"},_hoisted_12$1={class:"jg-td-c"},_hoisted_13$1={class:"jg-td-c"},_hoisted_14$1={class:"jg-td-c"},_hoisted_15={class:"jg-table-tools"},_hoisted_16={class:"jg-table-tool"},_hoisted_17=["onClick"],_hoisted_18={class:"jg-table-footer"},_hoisted_19={class:"jg-navigation"},_hoisted_20={class:"pagination"},_hoisted_21={class:"page-item"},_hoisted_22={key:0,class:"jg-loading"},_sfc_main$2={__name:"mange",emits:["create"],setup(m,{emit:p}){const d=p;let c=useRouter(),{currentRoute:{_rawValue:{params:{app_key:o}}}}=c;const i={NEXT:1,RESET:2};let h=getCurrentInstance(),n=reactive({params:{user_id:"",start:new Date(Date.now()-2*60*60*1e3),count:15},list:[],currentCount:0});function g(){let{start:_,count:t,user_id:s}=n.params;_=new Date(_).getTime(),C({start:_,app_key:o,count:t,user_id:s},i.RESET)}function B(){let{list:_,params:t}=n,{count:s,user_id:e}=t,a=_.length-1,l=_[a].timestamp;C({start:l,app_key:o,count:s,user_id:e},i.NEXT)}function C(_,t){if(n.isLoading)return;n.isLoading=!0;let{list:s}=n;Inspect.getConns(_).then(e=>{let{code:a,msg:l,data:u}=e;if(utils.extend(n,{isLoading:!1}),utils.isEqual(a,RESPONSE.SUCCESS)){let{logs:r}=u;r=utils.map(r,k=>{let{timestamp:E}=k;return k.connTimeName=N(E,"yyyy-MM-dd hh:mm:ss"),k}),utils.isEqual(t,i.NEXT)?s=s.concat(r):s=r,utils.extend(n,{list:s,currentCount:r.length})}else h.proxy.$toast({icon:"error",text:`Error: ${a} ${l}`})})}function v(_){d("create",{..._,count:20,start:_.timestamp,currentCount:0,isClose:!0,list:[],isLoading:!0})}function N(_,t="yyyy-MM-dd hh:mm"){return utils.formatTime(new Date(_).getTime(),t)}return(_,t)=>{const s=resolveComponent("VDatePicker");return openBlock(),createElementBlock("div",_hoisted_1$2,[createBaseVNode("div",_hoisted_2$2,[createBaseVNode("ul",_hoisted_3$2,[createBaseVNode("li",_hoisted_4$2,[createVNode(s,{modelValue:unref(n).params.start,"onUpdate:modelValue":t[0]||(t[0]=e=>unref(n).params.start=e),mode:"dateTime",is24hr:""},{default:withCtx(({togglePopover:e})=>[createBaseVNode("div",{class:"form-control jg-as-date-content",onClick:e},toDisplayString(N(unref(n).params.start)),9,_hoisted_5$2)]),_:1},8,["modelValue"])]),createBaseVNode("li",_hoisted_6$2,[withDirectives(createBaseVNode("input",{class:"form-control",type:"text","onUpdate:modelValue":t[1]||(t[1]=e=>unref(n).params.user_id=e),placeholder:"用户 ID",autocomplete:"off",onKeydown:withKeys(g,["enter"])},null,544),[[vModelText,unref(n).params.user_id]])]),createBaseVNode("li",{class:"jg-table-lf-item"},[createBaseVNode("div",{class:"jg-button jg-button-bg",onClick:g},"查询")])])]),createBaseVNode("div",_hoisted_7$2,[createBaseVNode("table",_hoisted_8$2,[t[2]||(t[2]=createBaseVNode("thead",null,[createBaseVNode("tr",null,[createBaseVNode("th",{class:"jg-td-c"},"连接时间"),createBaseVNode("th",{class:"jg-td-c"},"用户 ID"),createBaseVNode("th",{class:"jg-td-c"},"连接 ID"),createBaseVNode("th",{class:"jg-td-c"},"平台"),createBaseVNode("th",{class:"jg-td-c"},"客户端 IP"),createBaseVNode("th",{class:"jg-td-c"},"操作")])],-1)),createBaseVNode("tbody",null,[(openBlock(!0),createElementBlock(Fragment,null,renderList(unref(n).list,e=>(openBlock(),createElementBlock("tr",null,[createBaseVNode("td",_hoisted_9$1,toDisplayString(e.connTimeName),1),createBaseVNode("td",_hoisted_10$1,toDisplayString(e.user_id),1),createBaseVNode("td",_hoisted_11$1,toDisplayString(e.session),1),createBaseVNode("td",_hoisted_12$1,toDisplayString(e.platform),1),createBaseVNode("td",_hoisted_13$1,toDisplayString(e.client_ip),1),createBaseVNode("td",_hoisted_14$1,[createBaseVNode("ul",_hoisted_15,[createBaseVNode("li",_hoisted_16,[createBaseVNode("a",{class:"btn-link",href:"#",onClick:a=>v(e)},"查看详情",8,_hoisted_17)])])])]))),256))])])]),createBaseVNode("div",_hoisted_18,[createBaseVNode("nav",_hoisted_19,[createBaseVNode("ul",_hoisted_20,[createBaseVNode("li",_hoisted_21,[unref(n).currentCount>=unref(n).params.count?(openBlock(),createElementBlock("a",{key:0,class:"page-link",href:"#","aria-label":"Next",onClick:B},t[3]||(t[3]=[createBaseVNode("span",{"aria-hidden":"true"},"下一页",-1)]))):createCommentVNode("",!0)])])])]),unref(n).isLoading?(openBlock(),createElementBlock("div",_hoisted_22,t[4]||(t[4]=[createBaseVNode("div",{class:"loader-dot"},null,-1)]))):createCommentVNode("",!0)])}}};let Clipboard$1=class{constructor(){}createFake(p,d,c,o){let i=document.createElement("textarea");i.setAttribute("style","position: absolute;overflow: hidden;width: 0;height: 0;top: 0;left: 0;"),i.innerText=p,document.body.appendChild(i),i.select();try{document.execCommand(d),i.remove()}catch(h){o(h)}c()}copy(p,d,c){this.createFake(p,"copy",()=>{d&&d()},o=>{c&&c(o)})}};window.Clipboard=new Clipboard$1;const CLIPBOARD_ATTRIBUTE="clipboard",CLIPBOARD_ATTRIBUTE_SUCCESS="clipboard-success",CLIPBOARD_ATTRIBUTE_ERROR="clipboard-error",eventListenerClick=evt=>{window.Clipboard.copy(evt.target.getAttribute(CLIPBOARD_ATTRIBUTE),()=>{eval(evt.target.getAttribute(CLIPBOARD_ATTRIBUTE_SUCCESS))},err=>{eval(evt.target.getAttribute(CLIPBOARD_ATTRIBUTE_ERROR))})},observeElementClick=m=>{m.removeEventListener("click",eventListenerClick),m.addEventListener("click",eventListenerClick)};document.addEventListener("DOMSubtreeModified",m=>{if(!m.target||!m.target.querySelectorAll)return;let p=m.target.querySelectorAll("["+CLIPBOARD_ATTRIBUTE+"]");Array.from(p).forEach(d=>observeElementClick(d))});const Clipboard=window.Clipboard,_hoisted_1$1={class:"jg-signal-contanier"},_hoisted_2$1={class:"jg-signals-box",ref:"signals"},_hoisted_3$1={class:"jg-signal fadeinx"},_hoisted_4$1={class:"jg-signal-header"},_hoisted_5$1=["onClick"],_hoisted_6$1={class:"jg-signal-avatar"},_hoisted_7$1={class:"jg-signal-name"},_hoisted_8$1={class:"jg-signal-body"},_hoisted_9={class:"jg-signal-list"},_hoisted_10={key:0,class:"jg-table-footer"},_hoisted_11={class:"jg-navigation"},_hoisted_12={class:"pagination"},_hoisted_13={class:"page-item"},_hoisted_14={key:1,class:"jg-loading"},_sfc_main$1={__name:"signal",props:["conn"],emits:["save","hide","next","pre"],setup(m,{emit:p}){const d=getCurrentInstance(),c=m,o=p;reactive({}),watch(()=>c.conn.list,()=>{nextTick(()=>{let{signals:h}=d.refs;h&&(h.scrollTop=h.scrollHeight)})});function i(h){Clipboard.copy(h.timestamp,utils.noop,utils.noop),d.proxy.$toast({icon:"success",text:"【时间】已复制",duration:1500})}return(h,n)=>(openBlock(),createElementBlock("div",_hoisted_1$1,[createBaseVNode("ul",_hoisted_2$1,[(openBlock(!0),createElementBlock(Fragment,null,renderList(c.conn.list,g=>(openBlock(),createElementBlock("li",_hoisted_3$1,[createBaseVNode("div",_hoisted_4$1,[createBaseVNode("div",{class:"jg-signal-time",onClick:B=>i(g)},toDisplayString(g.timeName),9,_hoisted_5$1),createBaseVNode("div",_hoisted_6$1,[createBaseVNode("span",{class:normalizeClass(["jg-signal-icon jgicon",["jgicon-conn-"+g.type]])},null,2),createBaseVNode("span",_hoisted_7$1,toDisplayString(g.title),1)])]),createBaseVNode("div",_hoisted_8$1,[createBaseVNode("ul",_hoisted_9,[(openBlock(!0),createElementBlock(Fragment,null,renderList(g.infos,B=>(openBlock(),createElementBlock("li",{class:normalizeClass(["jg-signal-item",[B.cls?B.cls:""]])},toDisplayString(B.name)+" "+toDisplayString(B.value),3))),256))])])]))),256))],512),c.conn.currentCount>=c.conn.count?(openBlock(),createElementBlock("div",_hoisted_10,[createBaseVNode("nav",_hoisted_11,[createBaseVNode("ul",_hoisted_12,[createBaseVNode("li",_hoisted_13,[createBaseVNode("a",{class:"page-link",href:"#","aria-label":"Next",onClick:n[0]||(n[0]=g=>o("next",{item:c.conn}))},n[1]||(n[1]=[createBaseVNode("span",{"aria-hidden":"true"},"下一页",-1)]))])])])])):createCommentVNode("",!0),c.conn.isLoading?(openBlock(),createElementBlock("div",_hoisted_14,n[2]||(n[2]=[createBaseVNode("div",{class:"loader-dot"},null,-1)]))):createCommentVNode("",!0)]))}},_hoisted_1={class:"jg-tcon-container"},_hoisted_2={class:"jg-tcon-headers nav nav-tabs",ref:"navs"},_hoisted_3=["onClick"],_hoisted_4=["onClick"],_hoisted_5={class:"jg-tconn-title-userid"},_hoisted_6={class:"jg-tconn-title-time"},_hoisted_7=["onClick"],_hoisted_8={class:"jg-tcon-contents"},_sfc_main={__name:"connection",setup(m){let p=useRouter(),{currentRoute:{_rawValue:{params:{app_key:d}}}}=p,c=getCurrentInstance(),o=reactive({isShowEdit:!1,current:{session:"manger"},tabs:[{connTimeName:"连接信息",session:"manger",isActive:!0,isClose:!1,content:"123"}]});function i(t,s){utils.map(o.tabs,e=>(e.isActive=utils.isEqual(t.session,e.session),e)),o.current=t}function h(t){o.tabs.splice(t,1);let s=t-1,e=o.tabs[s];i(e)}function n(t){let s={...t},e=utils.find(o.tabs,a=>utils.isEqual(t.session,a.session));if(e>-1)return s=o.tabs[e],i(s);o.tabs.push(s),e=o.tabs.length-1,i(s),nextTick(()=>{let{navs:a}=c.refs;a&&(a.scrollLeft=a.scrollWidth)}),B({index:e,item:s,start:0},({index:a,list:l})=>{o.tabs[a].list=s.list.concat(l)})}function g({item:t}){let s=utils.find(o.tabs,u=>utils.isEqual(t.session,u.session)),e=o.tabs[s];e.isLoading=!0;let a=e.list.length-1,l=e.list[a];B({index:s,item:e,start:l.timestamp},({index:u,list:r})=>{e.list=e.list.concat(r)})}function B(t,s){let{item:e,index:a,start:l}=t,{session:u,user_id:r,count:k}=e;Inspect.getConn({start:l,session:u,user_id:r,count:k,app_key:d}).then(E=>{let{code:f,msg:$,data:V}=E;if(e.isLoading=!1,utils.isEqual(f,RESPONSE.SUCCESS)){let{logs:y}=V;e.currentCount=y.length,utils.isEqual(l,0)&&y.unshift({action:"connect",real_time:e.timestamp,user_id:e.user_id});let T=v(y);s({list:T,index:a})}else c.proxy.$toast({icon:"error",text:`Error: ${f} ${$}`})})}let C={connect:{type:SIGNAL_TYPE.CONNECTED,method:"connect"},qry:{type:SIGNAL_TYPE.USER},qry_ack:{type:SIGNAL_TYPE.REPLY,method:"qry_ack"},u_pub:{type:SIGNAL_TYPE.USER},u_pub_ack:{type:SIGNAL_TYPE.REPLY,method:"u_pub_ack"},s_pub:{type:SIGNAL_TYPE.SERVER},s_pub_ack:{type:SIGNAL_TYPE.REPLY,method:"s_pub_ack"},disconnect:{type:SIGNAL_TYPE.DISCONNECTED,method:"disconnect"}};function v(t){return console.log(t),utils.map(t,e=>{let{real_time:a,action:l}=e,u=utils.clone(e),r=N(u),k=C[l]||{};e=utils.extend(e,{...k,infos:r});let{method:E}=e,f=METHOD_MAP[E]||{};return e=utils.extend(e,f),e.timeName=_(a,"hh:mm:ss.S"),e})}function N(t){let{code:s}=t,e={name:"成功",value:"",cls:"success"};if(!utils.isUndefined(s)){let u=utils.find(IM_ERRORS,k=>utils.isEqual(k.code,s)),r=IM_ERRORS[u]||{code:s,msg:""};e={name:"失败",value:`: ${r.code} ${r.msg}`,cls:"warn"}}let a=["action","app_key","method","session","timestamp","code","real_time"];utils.forEach(a,u=>{delete t[u]});let l=[];return utils.forEach(t,(u,r)=>{l.push({name:r,value:`: ${u}`,order:r.charCodeAt(0)})}),l=utils.sort(l,(u,r)=>u.order>r.order),l.unshift(e),l}function _(t,s="yyyy-MM-dd hh:mm"){return utils.formatTime(new Date(t).getTime(),s)}return(t,s)=>(openBlock(),createElementBlock("div",_hoisted_1,[createBaseVNode("ul",_hoisted_2,[(openBlock(!0),createElementBlock(Fragment,null,renderList(unref(o).tabs,(e,a)=>(openBlock(),createElementBlock("li",{class:normalizeClass(["jg-tcon-nav-item nav-item fadeinx",[e.isClose?"":"jg-tcon-nav-item-ab"]])},[e.isClose?(openBlock(),createElementBlock("span",{key:0,class:"nav-close jgicon jgicon-close-c",onClick:l=>h(a)},null,8,_hoisted_3)):createCommentVNode("",!0),e.isClose?(openBlock(),createElementBlock("div",{key:1,class:normalizeClass(["nav-link jg-tcon-nav-item-link",[e.isActive?"active":""]]),href:"#",onClick:l=>i(e)},[createBaseVNode("span",_hoisted_5,toDisplayString(e.user_id),1),createBaseVNode("span",_hoisted_6,toDisplayString(e.connTimeName),1)],10,_hoisted_4)):(openBlock(),createElementBlock("div",{key:2,class:normalizeClass(["nav-link",[e.isActive?"active":""]]),href:"#",onClick:l=>i(e)},toDisplayString(e.connTimeName),11,_hoisted_7))],2))),256))],512),createBaseVNode("ul",_hoisted_8,[(openBlock(!0),createElementBlock(Fragment,null,renderList(unref(o).tabs,(e,a)=>(openBlock(),createElementBlock("li",{class:normalizeClass(["jg-tcon-content",[e.session==unref(o).current.session?"display-flex":"display-none"]])},[e.session=="manger"?(openBlock(),createBlock(_sfc_main$2,{key:0,onCreate:n})):(openBlock(),createBlock(_sfc_main$1,{key:1,conn:e,onNext:g},null,8,["conn"]))],2))),256))])]))}};export{_sfc_main as default};