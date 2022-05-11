package base62_impl

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func Test_Encode(t *testing.T) {
	type encodeVal struct {
		val []byte
	}
	type decodeData struct {
		input []byte
		want  encodeVal
	}
	testDataSet := []*decodeData{
		{
			input: []byte(""),
			want: encodeVal{
				val: nil,
			},
		},
		{
			input: []byte("1"),
			want: encodeVal{
				val: []byte("c8"),
			},
		}, {
			input: []byte("12"),
			want: encodeVal{
				val: []byte("c8N0"),
			},
		}, {
			input: []byte("123"),
			want: encodeVal{
				val: []byte("c8N5c"),
			},
		}, {
			input: []byte("1234"),
			want: encodeVal{
				val: []byte("c8N5cp0"),
			},
		}, {
			input: []byte("12345"),
			want: encodeVal{
				val: []byte("c8N5cp1F"),
			},
		}, {
			input: []byte("1234abcd"),
			want: encodeVal{
				val: []byte("c8N5cp51ohxIg"),
			},
		},
	}

	for _, td := range append(testDataSet) {
		require.Equal(t, td.want, encodeVal{
			val: b62encode(td.input),
		})
	}
}

func Test_Decode(t *testing.T) {
	type decodeVal struct {
		val []byte
		err error
	}
	type decodeData struct {
		input []byte
		want  decodeVal
	}
	testDataSet := []*decodeData{
		{
			input: []byte("c8"),
			want: decodeVal{
				val: []byte("1"),
			},
		}, {
			input: []byte("c8N0"),
			want: decodeVal{
				val: []byte("12"),
			},
		}, {
			input: []byte("c8N5c"),
			want: decodeVal{
				val: []byte("123"),
			},
		}, {
			input: []byte("c8N5cp0"),
			want: decodeVal{
				val: []byte("1234"),
			},
		}, {
			input: []byte("c8N5cp1F"),
			want: decodeVal{
				val: []byte("12345"),
			},
		}, {
			input: []byte("c8N5cp51ohxIg"),
			want: decodeVal{
				val: []byte("1234abcd"),
			},
		},
	}

	testDataSetForException := []*decodeData{
		{
			input: []byte("1"),
			want: decodeVal{
				err: InvalidInputLenErr(1),
			},
		}, {
			input: []byte("!2"),
			want: decodeVal{
				err: InvalidInputErr(0),
			},
		}, {
			input: []byte("2!"),
			want: decodeVal{
				err: InvalidInputErr(1),
			},
		}, {
			input: []byte("012345!"),
			want: decodeVal{
				err: InvalidInputErr(6),
			},
		},
	}
	for _, td := range append(testDataSet, testDataSetForException...) {
		val, err := b62decode(td.input)
		require.Equal(t, td.want, decodeVal{
			val: val,
			err: err,
		})
	}

}

func Test_All(t *testing.T) {

	type decodeData struct {
		input []byte
	}

	img1, err := ioutil.ReadFile("testdata/hbase.jpg")
	require.Nil(t, err)

	testDataSet := []*decodeData{
		{
			input: []byte("1234567890!@#$%^&*()_+"),
		}, {
			input: []byte("1234567890!@#$%^&*()-+_=" +
				"❤❥웃유♋☮✌☏☢☠✔☑♚▲♪✈✞÷↑↓◆◇⊙■□△▽¿─" +
				"Most encoding is the process of converting a string into binary" +
				"大多数编码都是由字符串转化成二进制的过程" +
				"Большинство кодировок - это процесс преобразования строки в двоичный" +
				"La plupart des encodages sont des processus de conversion de chaînes en binaires" +
				"अधिकतमा सङ्केतनम्" +
				"ការ​អ៊ិនកូដ​ភាព​ច្រើន​គឺ​ជា​ដំណើរការ​ការ​បម្លែង​ខ្សែអក្សរ​ទៅ​ជា​គោលពីរ" +
				"უფრო მეტი კოდირების პროცესია, რომელიც სტრიქონის ბინარიურად გადატანა" +
				"대부분의 인코딩은 문자열에서 2진법으로 바뀌는 과정이다"),
		}, {
			input: img1,
		},
	}
	for _, td := range testDataSet {
		output := b62encode(td.input)
		decoded, err := b62decode(output)

		require.Equal(t, err, nil)
		require.Equal(t, td.input, decoded)
	}
}
