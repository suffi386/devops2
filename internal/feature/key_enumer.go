// Code generated by "enumer -type Key -transform snake -trimprefix Key"; DO NOT EDIT.

package feature

import (
	"fmt"
	"strings"
)

const _KeyName = "unspecifiedlogin_default_orgtrigger_introspection_projectionslegacy_introspectionuser_schema"

var _KeyIndex = [...]uint8{0, 11, 28, 61, 81, 92}

const _KeyLowerName = "unspecifiedlogin_default_orgtrigger_introspection_projectionslegacy_introspectionuser_schema"

func (i Key) String() string {
	if i < 0 || i >= Key(len(_KeyIndex)-1) {
		return fmt.Sprintf("Key(%d)", i)
	}
	return _KeyName[_KeyIndex[i]:_KeyIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _KeyNoOp() {
	var x [1]struct{}
	_ = x[KeyUnspecified-(0)]
	_ = x[KeyLoginDefaultOrg-(1)]
	_ = x[KeyTriggerIntrospectionProjections-(2)]
	_ = x[KeyLegacyIntrospection-(3)]
	_ = x[KeyUserSchema-(4)]
}

var _KeyValues = []Key{KeyUnspecified, KeyLoginDefaultOrg, KeyTriggerIntrospectionProjections, KeyLegacyIntrospection, KeyUserSchema}

var _KeyNameToValueMap = map[string]Key{
	_KeyName[0:11]:       KeyUnspecified,
	_KeyLowerName[0:11]:  KeyUnspecified,
	_KeyName[11:28]:      KeyLoginDefaultOrg,
	_KeyLowerName[11:28]: KeyLoginDefaultOrg,
	_KeyName[28:61]:      KeyTriggerIntrospectionProjections,
	_KeyLowerName[28:61]: KeyTriggerIntrospectionProjections,
	_KeyName[61:81]:      KeyLegacyIntrospection,
	_KeyLowerName[61:81]: KeyLegacyIntrospection,
	_KeyName[81:92]:      KeyUserSchema,
	_KeyLowerName[81:92]: KeyUserSchema,
}

var _KeyNames = []string{
	_KeyName[0:11],
	_KeyName[11:28],
	_KeyName[28:61],
	_KeyName[61:81],
	_KeyName[81:92],
}

// KeyString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func KeyString(s string) (Key, error) {
	if val, ok := _KeyNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _KeyNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Key values", s)
}

// KeyValues returns all values of the enum
func KeyValues() []Key {
	return _KeyValues
}

// KeyStrings returns a slice of all String values of the enum
func KeyStrings() []string {
	strs := make([]string, len(_KeyNames))
	copy(strs, _KeyNames)
	return strs
}

// IsAKey returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Key) IsAKey() bool {
	for _, v := range _KeyValues {
		if i == v {
			return true
		}
	}
	return false
}