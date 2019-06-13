// Copyright 2016-2019 Jimmy zha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hessian

func init() {
	RegisterPOJO(&EnumConstantNotPresentException{})
	RegisterPOJO(&Class{})
}

type EnumConstantNotPresentException struct {
	SerialVersionUID     int64
	DetailMessage        string
	StackTrace           []StackTraceElement
	ConstantName         string
	EnumType             Class
	SuppressedExceptions []EnumConstantNotPresentException
	Cause                *EnumConstantNotPresentException
}

func NewEnumConstantNotPresentException(detailMessage string) *EnumConstantNotPresentException {
	return &EnumConstantNotPresentException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e EnumConstantNotPresentException) Error() string {
	return e.DetailMessage
}

func (EnumConstantNotPresentException) JavaClassName() string {
	return "java.lang.EnumConstantNotPresentException"
}