package manifest

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

const separator = "\n---\n"

func Decode(data string) ([]Resource, error) {
	docs := bytes.Split([]byte(data), []byte(separator))

	var resources []Resource
	for _, doc := range docs {
		r := &resource{}
		if err := yaml.Unmarshal(doc, r); err != nil {
			return nil, err
		}

		switch r.GetKind() {
		case KindPipeline:
			p := &Pipeline{}
			if err := yaml.Unmarshal(doc, p); err != nil {
				return nil, err
			}

			resources = append(resources, p)
		case KindSecret:
			s := &Secret{}
			if err := yaml.Unmarshal(doc, s); err != nil {
				return nil, err
			}

			resources = append(resources, s)
		default:
			resources = append(resources, r)
		}
	}

	return resources, nil
}

func Encode(resources []Resource) (string, error) {
	if len(resources) < 1 {
		return "", nil
	}

	buf := bytes.NewBuffer(nil)
	for idx, r := range resources {
		if idx != 0 {
			if _, err := buf.WriteString(separator); err != nil {
				return "", err
			}
		}

		resourceBytes, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}

		if _, err := buf.Write(resourceBytes); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}