import{O as v,d as r,P as f,r as O,V as P,s as I,c as D,a as s,t as g,u as l,I as k,p as w,F as G,b as B,h as L,w as q,i as $,C as F,o as S,a5 as x,g as W,a6 as T}from"./index-fb0c9453.js";let Y=c=>v(f.USER_DISABLE,{method:"GET",body:r.toJSON(c)}),J=c=>{let p=new URLSearchParams(c).toString(),d=`${f.ANALIYSIS_MESSAGE}?${p}`;return v(d,{method:"GET"})},z=c=>v(f.USER_DISABLE,{method:"GET",body:r.toJSON(c)}),H=c=>v(f.USER_DISABLE,{method:"GET",body:r.toJSON(c)});const M={getUserDAU:Y,getMessageChat:J,getGroupChat:z,getChatroomChat:H};function K(c){let m=c*24*60*60*1e3,p=r.formatTime(Date.now(),"yyyy-MM-dd"),d=new Date(`${p} 00:00:00`).getTime()-m,n=r.formatTime(d),y=new Date(n).getTime(),_=Date.now();return{start:y,end:_}}const R={getRangeDate:K},Q={class:"mb-4 jg-as-box"},X={class:"row jg-cb-row jg-as-header"},Z={class:"jg-bk-form"},tt={class:"row jg-asr-row"},et={class:"col-sm-4 jg-asr-col"},at={class:"jg-ars-num"},st={class:"jg-ars-percent"},nt={class:"col-sm-4 jg-asr-col"},ot={class:"jg-ars-num"},rt={class:"jg-ars-percent"},it={class:"jg-as-tools"},lt={class:"jg-as-tool"},ct=["onClick"],dt={class:"jg-as-date jgicon jgicon-date"},ut=["onClick"],gt={class:"row jg-as-body"},mt={class:"jg-bk-form",ref:"asuserchat"},yt={__name:"message",setup(c){const m=W();let p=$(),{currentRoute:{_rawValue:{params:{app_key:d}}}}=p,n=O({buttons:[{title:"7 天",name:7,isActive:!0},{title:"14 天",name:14,isActive:!1},{title:"30 天",name:30,isActive:!1}],range:{start:new Date(Date.now()-7*24*60*60*1e3),end:new Date},yestday:{p2p:{percent:0,isUp:!0,count:0},group:{percent:0,isUp:!0,count:0}}}),y=null;function _(e){if(console.log(e),!y){let{asuserchat:h}=m.refs;y=m.proxy.$echat.init(h)}let{privateData:t,groupData:o}=e,a=r.map(t.items,h=>r.formatTime(h.time_mark,"yyyy-MM-dd")).reverse(),i=C(t.items).reverse(),u=C(o.items).reverse();const b=["#5470C6","#EE6666"];let N={legend:{data:["日单聊下行消息量","日群聊下行消息量"]},tooltip:{trigger:"none",axisPointer:{type:"cross"}},xAxis:{type:"category",boundaryGap:!1,data:a},yAxis:{type:"value",axisPointer:{label:{formatter:function(h){return h.value.toFixed(0)}}}},series:[{name:"日单聊下行消息量",type:"line",smooth:!0,data:i,lineStyle:{color:b[1]}},{name:"日群聊下行消息量",type:"line",smooth:!0,lineStyle:{color:b[0]},data:u}]};y.setOption(N)}function C(e){return r.map(e,t=>t.count)}async function j(e){let{data:t}=await M.getMessageChat({...e,channel_type:x.PRIVATE}),{data:o}=await M.getMessageChat({...e,channel_type:x.GROUP}),a=(i,u)=>i.time_mark>u.time_mark;return t.items=r.sort(t.items||[],a),o.items=r.sort(o.items||[],a),{privateData:t,groupData:o}}function E(e){return r.formatTime(new Date(e).getTime(),"yyyy-MM-dd")}async function U(e){r.map(n.buttons,u=>(u.isActive=r.isEqual(u.name,e.name),u));let{start:t,end:o}=R.getRangeDate(e.name),a={app_key:d,start:t,end:o,stat_type:T.DOWN},i=await j(a);_(i)}P(()=>{let{start:e,end:t}=R.getRangeDate(7),o={app_key:d,start:e,end:t,stat_type:T.DOWN};j(o).then(a=>{V(a),_(a)})});function V(e){let{privateData:t,groupData:o}=e,a=A(t.items);r.extend(n.yestday.p2p,a);let i=A(o.items);r.extend(n.yestday.group,i)}function A(e){let t=e[0]||{count:0},o=e[1]||{count:1},a=(t.count-o.count)/o.count*100,i=!1;return a>=0&&(i=!0),{isUp:i,percent:Math.abs(a.toFixed(2)),count:r.numberWithCommas(t.count)}}return I(()=>n.range,async()=>{let{start:e,end:t}=n.range;e=new Date(e).getTime(),t=new Date(t).getTime();let o={app_key:d,start:e,end:t,stat_type:T.DOWN},a=await j(o);_(a)}),(e,t)=>{const o=F("VDatePicker");return S(),D("div",Q,[s("div",X,[s("div",Z,[s("div",tt,[s("div",et,[t[2]||(t[2]=s("span",{class:"jg-ars-memo"},"今日单聊下行消息量（条）",-1)),s("div",at,g(l(n).yestday.p2p.count),1),s("div",st,[t[1]||(t[1]=k(" 较前一日")),s("span",{class:w(["jgicon jg-ars-direction",[l(n).yestday.p2p.isUp?"jgicon-ac-up":"jgicon-ac-down"]])},g(l(n).yestday.p2p.percent)+"%",3)])]),s("div",nt,[t[4]||(t[4]=s("span",{class:"jg-ars-memo"},"今日群下行消息量（条）",-1)),s("div",ot,g(l(n).yestday.group.count),1),s("div",rt,[t[3]||(t[3]=k(" 较前一日")),s("span",{class:w(["jgicon jg-ars-direction",[l(n).yestday.group.isUp?"jgicon-ac-up":"jgicon-ac-down"]])},g(l(n).yestday.group.percent)+"%",3)])])])])]),s("div",it,[s("div",lt,[(S(!0),D(G,null,B(l(n).buttons,a=>(S(),D("div",{class:w(["jg-as-button",{"jg-as-button-active":a.isActive}]),onClick:i=>U(a)},g(a.title),11,ct))),256)),s("div",dt,[L(o,{modelValue:l(n).range,"onUpdate:modelValue":t[0]||(t[0]=a=>l(n).range=a),modelModifiers:{range:!0},class:"jg-as-date-picker"},{default:q(({togglePopover:a})=>[s("div",{class:"jg-as-date-content",onClick:a},g(E(l(n).range.start))+" 至 "+g(E(l(n).range.end)),9,ut)]),_:1},8,["modelValue"])])])]),s("div",gt,[s("div",mt,null,512)])])}}};export{yt as default};