import{r as b,s as E,d as c,o,c as r,a as i,t as k,j as N,v as V,u as t,x as B,F as g,b as y,i as U,A as S,p as C,y as $,z as x,k as j,g as F}from"./index-fb0c9453.js";const I={class:"jg-sw-form"},z={class:"jg-form-check form-switch"},A={class:"form-check-label"},D={__name:"input",props:["item"],emits:["save"],setup(h,{emit:d}){const e=h,p=d;let s=b({value:Number(e.item.value)});function l(){let{id:a}=e.item,{value:n}=s;p("save",{id:a,value:n})}return E(()=>e.item.value,a=>{c.extend(s,{value:Number(a)})}),(a,n)=>(o(),r("div",I,[i("div",z,[i("label",A,k(e.item.name),1),N(i("input",{class:"form-control",type:"number","onUpdate:modelValue":n[0]||(n[0]=m=>t(s).value=m)},null,512),[[V,t(s).value]]),i("div",{class:"jg-button",onClick:l},"保存")])]))}},L={class:"jg-sw-form"},M={class:"jg-form-check form-switch"},P={class:"form-check-label"},R=["checked"],H={__name:"switch",props:["item"],emits:["save"],setup(h,{emit:d}){const e=h,p=d;let s=b({value:Number(e.item.value)});function l(a){let{id:n}=e.item,m=a.target.checked,w=Number(m)+"";p("save",{id:n,value:w})}return E(()=>e.item.value,a=>{c.extend(s,{value:Number(a)})}),(a,n)=>(o(),r("div",L,[i("div",M,[i("label",P,k(e.item.name),1),i("input",{class:"form-check-input",type:"checkbox",checked:t(s).value,onChange:l},null,40,R)])]))}},W={class:"jg-sw-form"},Y={class:"jg-form-check form-switch"},G={class:"form-check-label"},J=["value"],K={__name:"select",props:["item"],emits:["save"],setup(h,{emit:d}){const e=h,p=d;let s=b({value:e.item.options[0].key});function l(){let{id:a}=e.item,{value:n}=s;p("save",{id:a,value:n})}return E(()=>e.item.value,a=>{c.extend(s,{value:a})}),(a,n)=>(o(),r("div",W,[i("div",Y,[i("label",G,k(e.item.name),1),N(i("select",{class:"form-select","onUpdate:modelValue":n[0]||(n[0]=m=>t(s).value=m)},[(o(!0),r(g,null,y(e.item.options,m=>(o(),r("option",{value:m.key},k(m.value),9,J))),256))],512),[[B,t(s).value]]),i("div",{class:"jg-button",onClick:l},"保存")])]))}},O={class:"md-4"},Q={class:"nav nav-underline-border",role:"tablist"},X=["onClick"],Z={class:"tab-content rounded-bottom"},ee={class:"row jg-sw-row"},te={class:"col-sm-4 jg-sw-col"},ae={__name:"switch",setup(h){let d=U(),{currentRoute:{_rawValue:{params:{app_key:e}}}}=d;const p=F();let s=[{type:"app",name:"应用相关",list:[{id:"token_effective_minute",type:"input",name:"Token 有效时长（小时）",value:0},{id:"kick_mode",type:"switch",name:"允许同设备多端登录",value:0}]},{type:"message",name:"消息相关",list:[{id:"is_hide_msg_before_join_group",type:"switch",name:"入群后获取之前的历史消息",value:0},{id:"his_msg_save_day",type:"select",name:"历史消息存储时长 (天)",value:"7",options:[{key:"7",value:"7 天"},{key:"360",value:"1 年"}]}]}],l=b({settings:s,current:s[0].type});function a(v){c.extend(l,{current:v.type})}function n(v){S.updateSetting({...v,app_key:e}).then(()=>{p.proxy.$toast({icon:"success",text:"保存成功"})})}function m(v,f){c.forEach(v,_=>{c.forEach(_.list,u=>{f(u)})})}function w(){let v=[];m(s,f=>{v.push(f.id)}),S.getSetting({app_key:e,config_keys:v}).then(({data:f})=>{let{configs:_}=f;m(l.settings,u=>{c.forEach(_,(T,q)=>{c.isEqual(u.id,q)&&(u.value=T)})})})}return w(),(v,f)=>(o(),r("div",O,[i("ul",Q,[(o(!0),r(g,null,y(t(l).settings,_=>(o(),r("li",{class:"nav-item sw-nav-item",onClick:u=>a(_)},[i("a",{class:C(["nav-link jgicon jgicon-free",{active:t(c).isEqual(t(l).current,_.type)}])},k(_.name),3)],8,X))),256))]),i("div",Z,[(o(!0),r(g,null,y(t(l).settings,_=>(o(),r("div",{class:C(["tab-pane p-3",{active:t(c).isEqual(t(l).current,_.type)}])},[i("div",ee,[(o(!0),r(g,null,y(_.list,u=>(o(),r("div",te,[t(c).isEqual(u.type,t($).INPUT)?(o(),x(D,{key:0,item:u,onSave:n},null,8,["item"])):j("",!0),t(c).isEqual(u.type,t($).SELECT)?(o(),x(K,{key:1,item:u,onSave:n},null,8,["item"])):j("",!0),t(c).isEqual(u.type,t($).SWITCH)?(o(),x(H,{key:2,item:u,onSave:n},null,8,["item"])):j("",!0)]))),256))])],2))),256))])]))}};export{ae as default};