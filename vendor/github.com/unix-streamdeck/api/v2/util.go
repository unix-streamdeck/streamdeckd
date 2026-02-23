package api


func CompareKeyConfigs(c1 KeyConfigV3, c2 KeyConfigV3) bool {
    return c1.Icon == c2.Icon &&
    c1.SwitchPage == c2.SwitchPage &&
    c1.Text == c2.Text &&
    c1.TextSize == c2.TextSize &&
    c1.TextAlignment == c2.TextAlignment &&
    c1.Keybind == c2.Keybind &&
    c1.Command == c2.Command &&
    c1.Brightness == c2.Brightness &&
    c1.Url == c2.Url &&
    c1.ObsCommand == c2.ObsCommand &&
    c1.IconHandler == c2.IconHandler &&
    c1.KeyHandler == c2.KeyHandler &&
    c1.Buff == c2.Buff &&
    c1.IconHandlerStruct == c2.IconHandlerStruct &&
    c1.KeyHandlerStruct == c2.KeyHandlerStruct
}


func CompareKeys(k1 KeyV3, k2 KeyV3) bool {
    for key, configk1 := range k1.Application {
        configk2, ok := k2.Application[key]
        if !ok {
            return false
        }
        if !CompareKeyConfigs(*configk1, *configk2) {
            return false
        }
    }
    for key, configk2 := range k2.Application {
        configk1, ok := k1.Application[key]
        if !ok {
            return false
        }
        if !CompareKeyConfigs(*configk1, *configk2) {
            return false
        }
    }
    return true
}