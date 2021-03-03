export default class UEInfoDetail {
  
  ueInfoDetail = {
      ocfInfo:{},
      smfInfo:{},
      pcfInfo:{}
  }

  constructor(info) {
     this.ocfInfo = info.ocfInfo;
     this.smfInfo = info.smfInfo;
     this.pcfInfo = info.pcfInfo;
  }
}