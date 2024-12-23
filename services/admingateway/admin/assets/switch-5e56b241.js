import{r as b,s as E,d as o,o as c,c as r,a as i,t as y,j as N,v as V,u as s,x as B,F as k,b as g,i as U,A as S,p as C,y as x,z as $,k as j,g as F}from"./index-a6509667.js";const I={class:"jg-sw-form"},z={class:"jg-form-check form-switch"},A={class:"form-check-label"},D={__name:"input",props:["item"],emits:["save"],setup(h,{emit:p}){const t=h,d=p;let a=b({value:Number(t.item.value)});function l(){let{id:n}=t.item,{value:e}=a;d("save",{id:n,value:e})}return E(()=>t.item.value,n=>{let e=Number(n);e>0&&o.extend(a,{value:e})}),(n,e)=>(c(),r("div",I,[i("div",z,[i("label",A,y(t.item.name),1),N(i("input",{class:"form-control",type:"number","onUpdate:modelValue":e[0]||(e[0]=m=>s(a).value=m)},null,512),[[V,s(a).value]]),i("div",{class:"jg-button",onClick:l},"保存")])]))}},L={class:"jg-sw-form"},M={class:"jg-form-check form-switch"},P={class:"form-check-label"},R=["checked"],H={__name:"switch",props:["item"],emits:["save"],setup(h,{emit:p}){const t=h,d=p;let a=b({value:Number(t.item.value)});function l(n){let{id:e}=t.item,m=n.target.checked,w=Number(m)+"";d("save",{id:e,value:w})}return E(()=>t.item.value,n=>{o.extend(a,{value:Number(n)})}),(n,e)=>(c(),r("div",L,[i("div",M,[i("label",P,y(t.item.name),1),i("input",{class:"form-check-input",type:"checkbox",checked:s(a).value,onChange:l},null,40,R)])]))}},W={class:"jg-sw-form"},Y={class:"jg-form-check form-switch"},G={class:"form-check-label"},J=["value"],K={__name:"select",props:["item"],emits:["save"],setup(h,{emit:p}){const t=h,d=p;let a=b({value:t.item.options[0].key});function l(){let{id:n}=t.item,{value:e}=a;d("save",{id:n,value:e})}return E(()=>t.item.value,n=>{o.extend(a,{value:n})}),(n,e)=>(c(),r("div",W,[i("div",Y,[i("label",G,y(t.item.name),1),N(i("select",{class:"form-select","onUpdate:modelValue":e[0]||(e[0]=m=>s(a).value=m)},[(c(!0),r(k,null,g(t.item.options,m=>(c(),r("option",{value:m.key},y(m.value),9,J))),256))],512),[[B,s(a).value]]),i("div",{class:"jg-button",onClick:l},"保存")])]))}},O={class:"md-4"},Q={class:"nav nav-underline-border",role:"tablist"},X=["onClick"],Z={class:"tab-content rounded-bottom"},ee={class:"row jg-sw-row"},te={class:"col-sm-4 jg-sw-col"},ae={__name:"switch",setup(h){let p=U(),{currentRoute:{_rawValue:{params:{app_key:t}}}}=p;const d=F();let a=[{type:"app",name:"应用相关",list:[{id:"token_effective_minute",type:"input",name:"Token 有效时长（小时）",value:0},{id:"kick_mode",type:"switch",name:"允许同设备多端登录",value:0}]},{type:"message",name:"消息相关",list:[{id:"is_hide_msg_before_join_group",type:"switch",name:"入群后获取之前的历史消息",value:0},{id:"not_check_grp_member",type:"switch",name:"不在群组是否可以获取群消息",value:0},{id:"his_msg_save_day",type:"select",name:"历史消息存储时长 (天)",value:"7",options:[{key:"7",value:"7 天"},{key:"360",value:"1 年"}]}]},{type:"group",name:"群组相关",list:[{id:"max_grp_member_count",type:"input",name:"群人数上限",value:1e3}]},{type:"chatroom",name:"聊天室相关",list:[{id:"chrm_msg_cache_max_count",type:"input",name:"单个聊天室消息桶大小",value:50},{id:"chrm_att_max_count",type:"input",name:"单个聊天室属性数量",value:100},{id:"chrm_event_ntf",type:"switch",name:"是否开启聊天室事件通知",value:!1},{id:"chrm_event_cache_max_count",type:"input",name:"单个聊天室事件桶大小",value:50}]}],l=b({settings:a,current:a[0].type});function n(v){o.extend(l,{current:v.type})}function e(v){S.updateSetting({...v,app_key:t}).then(()=>{d.proxy.$toast({icon:"success",text:"保存成功"})})}function m(v,f){o.forEach(v,_=>{o.forEach(_.list,u=>{f(u)})})}function w(){let v=[];m(a,f=>{v.push(f.id)}),S.getSetting({app_key:t,config_keys:v}).then(({data:f})=>{let{configs:_}=f;m(l.settings,u=>{o.forEach(_,(T,q)=>{o.isEqual(u.id,q)&&(u.value=T)})})})}return w(),(v,f)=>(c(),r("div",O,[i("ul",Q,[(c(!0),r(k,null,g(s(l).settings,_=>(c(),r("li",{class:"nav-item sw-nav-item",onClick:u=>n(_)},[i("a",{class:C(["nav-link jgicon jgicon-free",{active:s(o).isEqual(s(l).current,_.type)}])},y(_.name),3)],8,X))),256))]),i("div",Z,[(c(!0),r(k,null,g(s(l).settings,_=>(c(),r("div",{class:C(["tab-pane p-3",{active:s(o).isEqual(s(l).current,_.type)}])},[i("div",ee,[(c(!0),r(k,null,g(_.list,u=>(c(),r("div",te,[s(o).isEqual(u.type,s(x).INPUT)?(c(),$(D,{key:0,item:u,onSave:e},null,8,["item"])):j("",!0),s(o).isEqual(u.type,s(x).SELECT)?(c(),$(K,{key:1,item:u,onSave:e},null,8,["item"])):j("",!0),s(o).isEqual(u.type,s(x).SWITCH)?(c(),$(H,{key:2,item:u,onSave:e},null,8,["item"])):j("",!0)]))),256))])],2))),256))])]))}};export{ae as default};
