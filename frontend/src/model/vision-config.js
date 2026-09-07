import $api from "common/api";
import Model from "./model";

export class VisionConfig extends Model {
  getDefaults() {
    return {
      models: [],
      thresholds: {},
    };
  }

  load() {
    return $api.get("config/vision").then((response) => {
      return Promise.resolve(this.setValues(response.data));
    });
  }

  save() {
    return $api.post("config/vision", this.getValues(true)).then((response) => Promise.resolve(this.setValues(response.data)));
  }
}

export default VisionConfig;
