using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using CsvHelper;
using Services;
using CsvLibrary;

namespace awstaghelper
{
    class Program
    {
        static void Main(string[] args)
        {
            Csv csv = new Csv();
            Ec2 ec2 = new Ec2(csv);
            ec2.ReadCsv();
        }
    }
}